package simulating

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/amenzhinsky/iothub/common"
	"github.com/amenzhinsky/iothub/iotdevice"
	iotmqtt "github.com/amenzhinsky/iothub/iotdevice/transport/mqtt"
	"github.com/amenzhinsky/iothub/logger"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type (
	// device represents the IoT Central device being simulated.
	device struct {
		deviceID             string                   // unique id of the device.
		model                *models.DeviceModel      // model of the device.
		target               *models.SimulationTarget // target application of the device.
		connectionString     string                   // IoT Hub connectionString of the device.
		isConnected          bool                     // is the device connected.
		isConnecting         bool                     // is the device connecting now.
		telemetrySentTime    time.Time                // last time telemetry was sent from this device.
		sendingTelemetry     bool                     // is the device sending telemetry now.
		sendingReportedProps bool                     // is the device sending reported properties now.
		iotHubClient         *iotdevice.Client        // IoT Hub connection MQTT client.
		twinSub              *iotdevice.TwinStateSub  // subscription to listen for twin updates.
		c2dSub               *iotdevice.EventSub      // subscription to listen for c2d commands
		dataGenerator        *DataGenerator           // data generator used to generate telemetry and reported property updates.
		retryCount           int                      // number of retries for sending telemetry
	}

	// deviceCollection represents collection of devices used in device groups.
	deviceCollection struct {
		devices []*device // devices in the collection.
	}

	// deviceSimulator is responsible for simulating device behaviors such as sending telemetry messages, reported properties, acknowledging twin updates and commands.
	deviceSimulator struct {
		cancel                context.CancelFunc         // cancel function to invoke when the device simulator is being stopped.
		context               context.Context            // the context of the device simulator.
		simulation            *models.Simulation         // the simulation that is driving this simulator.
		config                *Config                    // starling level simulation configuration.
		telemetryRequests     chan *telemetryRequest     // input channel used for queuing up telemetry requests.
		reportedPropsRequests chan *reportedPropsRequest // input channel used for queuing up reported property requests.
		provisioner           *DeviceProvisioner         // provisioner to provision devices using DPS
		provisionThrottle     chan int                   // channel to apply device provisioning rate throttle
	}

	// telemetryRequest represents the request to send telemetry by the device simulator.
	telemetryRequest struct {
		device  *device         // device the device which sends telemetry.
		context context.Context // context of the telemetry request.
	}

	// telemetryMessage represents the telemetry message sent from the device.
	telemetryMessage struct {
		body               []byte            // body of the telemetry message.
		interfaceId        string            // interface id of the component that is sending telemetry.
		connectionDeviceID string            // id of the device sending telemetry.
		connectionModuleID string            // edge module that is sending telemetry.
		contentEncoding    string            // encoding of the message content.
		contentType        string            // content type of the body.
		correlationID      string            // correlation is to be used from IoT Hub downstream systems.
		messageID          string            // unique identifier used to correlate two-way communication.
		creationTimeUtc    time.Time         // time when the device generated telemetry.
		properties         map[string]string // telemetry message headers sent by the device.
		dataPointCount     int               // number of data points sent in the message
	}

	// telemetryBatch represents a batch of telemetry messages.
	telemetryBatch struct {
		messages []*telemetryMessage // a collection of telemetry messages distributed since the last telemetry send
	}

	// reportedPropsRequest represents a request to send reported properties update from the device.
	reportedPropsRequest struct {
		device  *device         // device the device which sends reported properties update.
		context context.Context // context of the reported properties update request.
	}
)

var (
	useMock bool = false // use mock (local random sleeps) instead of real IoT Hub interactions
)

// newDeviceSimulator create a new device simulator
func newDeviceSimulator(ctx context.Context, config *Config, simulation *models.Simulation) *deviceSimulator {
	deviceSimContext, cancel := context.WithCancel(ctx)
	return &deviceSimulator{
		cancel:                cancel,
		context:               deviceSimContext,
		config:                config,
		simulation:            simulation,
		telemetryRequests:     make(chan *telemetryRequest, config.MaxConcurrentConnections),     // only process so many concurrent telemetry requests at a time
		reportedPropsRequests: make(chan *reportedPropsRequest, config.MaxConcurrentConnections), // only process so many concurrent reported property update send requests at a time
		provisioner:           NewProvisioner(ctx, config),
		provisionThrottle:     make(chan int, config.MaxConcurrentRegistrations), // only allow so many DPS registrations at a time
	}
}

// start starts the telemetry and reported property update pumps in this device simulator
func (s *deviceSimulator) start(totalDevices int) {
	// start telemetry pump if it is enabled in config file
	if s.config.EnableTelemetry {
		// take the min(maxConnections, totalDevices)
		maxConnections := totalDevices
		if maxConnections > s.config.MaxConcurrentConnections {
			maxConnections = s.config.MaxConcurrentConnections
		}

		// create parallel telemetry request processors
		for i := 1; i <= maxConnections; i++ {
			go func(pumpId int) {
				for {
					select {
					case <-s.context.Done():
						log.Trace().Int("pumpId", pumpId).Msg("device simulation pump stopped")
						return
					case telemetryReq := <-s.telemetryRequests:
						s.sendTelemetry(telemetryReq)
					}
				}
			}(i)
		}
		log.Debug().Msg("device simulation telemetry consumer pump started")
	}

	// start reported properties pump if it is enabled in config file
	if s.config.EnableReportedProps {
		// create parallel reported props update request processors
		for i := 1; i <= s.config.MaxConcurrentTwinUpdates; i++ {
			go func(pumpId int) {
				for {
					select {
					case <-s.context.Done():
						log.Trace().Int("pumpId", pumpId).Msg("device simulation pump stopped")
						return
					case reportedPropsReq := <-s.reportedPropsRequests:
						s.sendReportedProps(reportedPropsReq)
					}
				}
			}(i)
		}
		log.Debug().Msg("device simulation reported properties consumer pump started")
	}
}

// stop stops the device simulator
func (s *deviceSimulator) stop() {
	// Stop all activity
	s.cancel()
}

// sendTelemetry sends a telemetry batch from the device
func (s *deviceSimulator) sendTelemetry(req *telemetryRequest) {
	// if the device is in the middle of sending a telemetry, skip this request
	if req.device.sendingTelemetry {
		log.Trace().
			Str("deviceID", req.device.deviceID).
			Msg("skipping telemetry as it is already sending one")
		telemetryBatchSkippedTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(1)
		return
	}

	req.device.sendingTelemetry = true

	// if there are too many retries, device might have disconnected or failed over; provision it again
	failureDetected := false
	if req.device.retryCount > 1 {
		// && req.device.isConnected == true {
		s.disconnectDevice(req.device)
		req.device.connectionString = ""
		// clear device from cache
		storing.TargetDevices.Delete(req.device.target.ID, req.device.deviceID)
		log.Debug().Str("deviceID", req.device.deviceID).Int("retryCount", req.device.retryCount).Msg("device might have been moved so will be re-provisioned")
		failureDetected = true
	}

	// make sure that the device is connected
	if req.device.isConnected == false {
		if s.connectDevice(req.device) == false {
			req.device.sendingTelemetry = false
			return
		}

		// device failed over successfully
		if failureDetected {
			deviceFailoverTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Inc()
		}
	}

	// generate a batch of telemetry messages
	batch := s.getNextTelemetryBatch(req.device)
	start := time.Now()

	// send all messages in a batch in parallel.
	wg := sync.WaitGroup{}
	for _, msg := range batch.messages {
		wg.Add(1)
		go s.sendTelemetryMessage(msg, req, &wg)
	}

	// wait till all messages in the batch are sent.
	wg.Wait()

	now := time.Now()
	req.device.telemetrySentTime = now
	latency := float64(now.UnixNano()-start.UnixNano()) / float64(time.Second)

	log.Trace().
		Str("deviceID", req.device.deviceID).
		Int("batchSize", len(batch.messages)).
		Float64("latency", latency).
		Int("numGoroutines", runtime.NumGoroutine()).
		Msg("sent telemetry")

	telemetryBatchSendLatency.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Observe(latency)
	telemetryBatchSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(1)

	// disconnect device based on the disconnect behavior
	if s.simulation.DisconnectBehavior == models.DeviceDisconnectAfterTelemetrySend {
		s.disconnectDevice(req.device)
	}

	req.device.sendingTelemetry = false
}

// sendTelemetryMessage sends a telemetry message from the device to IoT hub
func (s *deviceSimulator) sendTelemetryMessage(msg *telemetryMessage, req *telemetryRequest, wg *sync.WaitGroup) bool {
	defer wg.Done()

	start := time.Now()
	var err error
	if useMock {
		randSleep(s.context, 500, 5000)
	} else {
		// send telemetry to IoT Central
		timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.TelemetryTimeout))
		err = req.device.iotHubClient.SendEvent(timeoutCtx, msg.body,
			iotdevice.WithSendCorrelationID(msg.correlationID),
			iotdevice.WithSendMessageID(msg.messageID),
			iotdevice.WithSendProperties(map[string]string{
				"iothub-creation-time-utc":    msg.creationTimeUtc.Format(time.RFC3339),
				"iothub-connection-device-id": msg.connectionDeviceID,
				"iothub-interface-id":         msg.interfaceId,
			}))
	}
	if err != nil {
		log.Error().
			Str("deviceID", req.device.deviceID).
			Err(err).
			Msg("error sending telemetry to hub")
		telemetryMessageFailureTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID, s.getErrorType(err)).Add(1)
		req.device.retryCount++
		return false
	} else {
		req.device.retryCount = 0
		telemetryMessageSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(1)
		latency := float64(time.Now().UnixNano()-start.UnixNano()) / float64(time.Second)
		telemetryMessageSendLatency.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Observe(latency)
		telemetrySentBytes.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(float64(len(msg.body)))
		telemetryDataPointsSentTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(float64(msg.dataPointCount))
	}
	return true
}

// sendReportedProps send reported properties update from the device
func (s *deviceSimulator) sendReportedProps(req *reportedPropsRequest) {
	// if the device is in the middle of sending a reported property update, skip this request
	if req.device.sendingReportedProps {
		log.Trace().
			Str("deviceID", req.device.deviceID).
			Msg("skipping reported properties as it is already sending one")
		reportedPropsSkippedTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(1)
		return
	}

	req.device.sendingReportedProps = true

	// make sure that the device is connected
	if req.device.isConnected == false {
		if s.connectDevice(req.device) == false {
			req.device.sendingReportedProps = false
			return
		}
	}

	//desired, reported, _ := req.device.iotHubClient.RetrieveTwinState(s.context)
	//log.Debug().Str("deviceID", req.device.deviceID).Msg(fmt.Sprintf("current desired: %v, reported: %v", desired, reported))

	// generate reported properties
	start := time.Now()
	reportedProps, err := req.device.dataGenerator.GenerateReportedProperties(req.device)
	if err != nil {
		reportedPropsFailureTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID, s.getErrorType(err)).Add(1)
		log.Debug().Err(err).Str("deviceID", req.device.deviceID).Msg("error generating reported property update")
	}
	log.Trace().Str("deviceID", req.device.deviceID).Msg(fmt.Sprintf("about to update reported props: %v", reportedProps))

	// send the reported properties to IoT Central
	timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.TwinUpdateTimeout))
	_, err = req.device.iotHubClient.UpdateTwinState(timeoutCtx, reportedProps)
	if err != nil {
		reportedPropsFailureTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID, s.getErrorType(err)).Add(1)
		log.Debug().Err(err).Str("deviceID", req.device.deviceID).Msg("error sending reported properties update")
		req.device.retryCount++
	} else {
		now := time.Now()
		latency := float64(now.UnixNano()-start.UnixNano()) / float64(time.Second)
		reportedPropsSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Add(1)
		reportedPropsSendLatency.WithLabelValues(s.simulation.ID, s.simulation.TargetID, req.device.model.ID).Observe(latency)
		log.Trace().
			Str("deviceID", req.device.deviceID).
			Float64("latency", latency).
			Int("numGoroutines", runtime.NumGoroutine()).
			Msg("sent reported properties")
		req.device.retryCount = 0
	}

	req.device.sendingReportedProps = false
}

// provisionDevice provision the device in Central
func (s *deviceSimulator) provisionDevice(device *device, useCache bool) bool {
	// see if the cache contains previously provisioned device
	if useCache {
		td, _ := storing.TargetDevices.Get(device.target.ID, device.deviceID)
		if td != nil {
			log.Trace().Str("deviceId", device.deviceID).Msg("found device in cache")
			device.connectionString = td.ConnectionString
			return true
		}
	}

	// apply provisioning throttle
	s.provisionThrottle <- 0

	// provision the device for the first time
	req := &ProvisioningRequest{
		DeviceID:   device.deviceID,
		Context:    s.context,
		Target:     device.target,
		Simulation: s.simulation,
		Model:      device.model,
	}
	result := s.provisioner.Provision(req)
	if result == nil {
		// remove provisioning throttle
		<-s.provisionThrottle
		return false
	}
	device.connectionString = result.ConnectionString

	// cache the device for future use
	newDevice := models.SimulationTargetDevice{
		TargetID:         req.Target.ID,
		DeviceID:         req.DeviceID,
		ConnectionString: result.ConnectionString,
	}
	if err := storing.TargetDevices.Set(&newDevice); err != nil {
		// remove provisioning throttle
		<-s.provisionThrottle
		log.Error().Err(err).Str("deviceID", device.deviceID).Msg("error caching device connection string")
		return false
	}
	log.Trace().Str("deviceId", device.deviceID).Str("connectionString", result.ConnectionString).Msg("saved device in cache")

	// remove provisioning throttle
	<-s.provisionThrottle
	return true
}

// connectDevice provisions and connects a given device
func (s *deviceSimulator) connectDevice(device *device) bool {
	// provision the device for the first time
	if len(device.connectionString) == 0 {
		if s.provisionDevice(device, true) == false {
			return false
		}
	}

	// if the device is in the middle of connecting, ignore this request
	if device.isConnecting {
		return false
	}

	device.isConnecting = true

	connectTimer := prometheus.NewTimer(deviceConnectLatency.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID))
	defer connectTimer.ObserveDuration()

	hub := getHubName(device.connectionString)

	if useMock {
		randSleep(s.context, 500, 1000)
	} else {
		var err error

		// connect the device to IoT Central
		device.iotHubClient, err = iotdevice.NewFromConnectionString(iotmqtt.New(), device.connectionString,
			iotdevice.WithLogger(logger.New(logger.LevelDebug, func(lvl logger.Level, s string) {
				log.Trace().Msg(s)
			})))
		if err != nil {
			device.isConnecting = false
			log.Error().Err(err).Str("deviceID", device.deviceID).Str("connectionString", device.connectionString).Msg("error parsing connection string")
			return false
		}

		log.Trace().Str("deviceID", device.deviceID).Str("connectionString", device.connectionString).Msg("trying to connect to iothub")
		timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.ConnectionTimeout))
		if err = device.iotHubClient.Connect(timeoutCtx); err != nil {
			device.isConnecting = false
			log.Error().Err(err).Str("deviceID", device.deviceID).Msg("error connecting to IoT Hub")

			// device might have moved to a different hub, provision and connect to hub again
			errMsg := strings.ToLower(err.Error())
			if errMsg == "not authorized" || errMsg == "server unavailable" || strings.Contains(errMsg, "network error") {
				log.Trace().Str("deviceID", device.deviceID).Msg("detected hub fail over, re-provisioning device")

				if s.provisionDevice(device, false) == false {
					return false
				}

				// close existing hub connections
				_ = device.iotHubClient.Close()
				device.iotHubClient = nil

				hub = getHubName(device.connectionString)
				device.iotHubClient, _ = iotdevice.NewFromConnectionString(iotmqtt.New(), device.connectionString,
					iotdevice.WithLogger(logger.New(logger.LevelDebug, func(lvl logger.Level, s string) {
						log.Trace().Msg(s)
					})))
				timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.ConnectionTimeout))
				if err = device.iotHubClient.Connect(timeoutCtx); err != nil {
					log.Error().Err(err).Str("deviceID", device.deviceID).Str("connectionString", device.connectionString).Msg("error connecting to IoT Hub")
					return false
				}
				deviceFailoverTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID).Inc()
				log.Debug().Str("deviceID", device.deviceID).Msg("detected hub fail over, reconnected to IoT Hub")
			} else {
				return false
			}
		}
		log.Trace().Err(err).Str("deviceID", device.deviceID).Msg("device connected to IoT Hub")

		// register for twin updates
		if s.config.EnableTwinUpdateAcks {
			if s.subscribeTwinUpdates(device) == false {
				device.isConnecting = false
				return false
			}
		}

		// register for c2d commands
		if s.config.EnableCommandAcks {
			if s.subscribeCommands(device) == false {
				device.isConnecting = false
				return false
			}
		}
	}

	device.isConnected = true
	device.isConnecting = false
	connectedDeviceGauge.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID, hub).Inc()

	return true
}

// disconnectDevice disconnects a given device from IoT Central
func (s *deviceSimulator) disconnectDevice(device *device) bool {

	hub := getHubName(device.connectionString)

	if useMock {
		randSleep(s.context, 500, 1000)
	} else {
		if device.iotHubClient != nil {
			// unregister for twin updates
			if s.config.EnableTwinUpdateAcks {
				s.unsubscribeTwinUpdates(device)
			}

			// unregister from c2d commands and direct methods
			if s.config.EnableCommandAcks {
				s.unsubscribeCommands(device)
			}

			_ = device.iotHubClient.Close()
			device.iotHubClient = nil
		}
	}
	log.Trace().Str("deviceID", device.deviceID).Msg("disconnected device from IoT Hub")

	// do not reset connection string
	// we reuse the connection string until we get a failure

	if device.isConnected {
		connectedDeviceGauge.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID, hub).Dec()
	}
	device.isConnected = false

	return true
}

// subscribeTwinUpdates creates subscription to monitor twin update (desired property) requests fpr a given device
func (s *deviceSimulator) subscribeTwinUpdates(device *device) bool {
	var err error
	timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.TwinUpdateTimeout))
	device.twinSub, err = device.iotHubClient.SubscribeTwinUpdates(timeoutCtx)
	if err != nil {
		// TODO: add retry
		log.Err(err).Str("deviceID", device.deviceID).Msg("twin update subscription failed")
		return false
	}

	go func() {
		for {
			select {
			case <-s.context.Done():
				log.Trace().Str("deviceID", device.deviceID).Msg("device twin subscription stopped")
				return
			case desiredTwin := <-device.twinSub.C():
				dt, _ := json.Marshal(desiredTwin)
				log.Trace().Str("deviceID", device.deviceID).
					Str("desiredTwin", fmt.Sprintf("%s", dt)).
					Msg("got twin update")

				// acknowledge twin update by echoing reported properties
				reportedTwin := device.dataGenerator.GenerateTwinUpdateAck(desiredTwin)
				start := time.Now()
				timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.TwinUpdateTimeout))
				_, err := device.iotHubClient.UpdateTwinState(timeoutCtx, reportedTwin)
				end := time.Now()
				latency := float64(end.UnixNano()-start.UnixNano()) / float64(time.Second)

				if err != nil {
					twinUpdateFailureTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID, s.getErrorType(err)).Add(1)
					log.Err(err).Str("deviceID", device.deviceID).Msg("twin update failed")
				} else {
					twinUpdateSendLatency.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID).Observe(latency)
					twinUpdateSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID).Add(1)
					rt, _ := json.Marshal(reportedTwin)
					log.Trace().Str("deviceID", device.deviceID).
						//Int("reportedVersion", reportedVersion).
						Str("reportedProperties", fmt.Sprintf("%s", rt)).
						Msg("acknowledged twin update")
				}
			}
		}
	}()

	return true
}

// unsubscribeTwinUpdates unsubscribe from twin updates for a given device
func (s *deviceSimulator) unsubscribeTwinUpdates(device *device) bool {
	if device.twinSub != nil {
		device.iotHubClient.UnsubscribeTwinUpdates(device.twinSub)
	}
	return true
}

// subscribeCommands subscribe for c2d command requests from IoT Central to the device
func (s *deviceSimulator) subscribeCommands(device *device) bool {
	// register for (Sync) Direct Methods
	hasAsyncCommands := false
	for _, component := range device.dataGenerator.CapabilityModel.Components {
		for _, command := range component.Commands {
			if command.IsSync {
				timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.CommandTimeout))
				err := device.iotHubClient.RegisterMethod(timeoutCtx, command.Name, func(p map[string]interface{}) (map[string]interface{}, error) {
					// acknowledge the c2d command by a reply
					// TODO: need to figure out how to respond with proper return types based on the DCM
					resp := make(map[string]interface{})
					commandsSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID).Add(1)
					log.Trace().Str("deviceID", device.deviceID).Str("Method", command.Name).Msg("direct method acknowledged")
					return resp, nil
				})

				if err != nil {
					log.Err(err).Str("deviceID", device.deviceID).Str("Method", command.Name).Msg("failed to register direct method")
					return false
				}
			} else {
				hasAsyncCommands = true
			}
		}
	}

	// register for C2D (Async) Commands
	if hasAsyncCommands {
		var err error
		timeoutCtx, _ := context.WithTimeout(s.context, time.Millisecond*time.Duration(s.config.CommandTimeout))
		device.c2dSub, err = device.iotHubClient.SubscribeEvents(timeoutCtx)
		if err != nil {
			log.Err(err).Str("deviceID", device.deviceID).Msg("c2d command subscription failed")
			return false
		}
		go func() {
			for {
				select {
				case <-s.context.Done():
					log.Trace().Str("deviceID", device.deviceID).Msg("c2d subscription stopped")
					return
				case msg := <-device.c2dSub.C():
					if msg != nil {
						log.Trace().Str("msg", string(msg.Properties["method-name"])).Str("msg", fmt.Sprintf("%v", msg)).Msg("received c2d command")
						commandsSuccessTotal.WithLabelValues(s.simulation.ID, s.simulation.TargetID, device.model.ID).Add(1)

						// send ack to c2d command
						//s.sendC2DAck(device, msg)
					}
				}
			}
		}()
	}
	return true
}

func (s *deviceSimulator) sendC2DAck(device *device, msg *common.Message) {
	//log.Debug().Str("msg", string(msg.Properties["method-name"])).Msg("received c2d command")
}

// unsubscribeCommands unsubscribe from c2d command requests for a given device
func (s *deviceSimulator) unsubscribeCommands(device *device) bool {
	for _, component := range device.dataGenerator.CapabilityModel.Components {
		for _, command := range component.Commands {
			device.iotHubClient.UnregisterMethod(command.Name)
		}
	}
	return true
}

// getNextTelemetryBatch creates a batch of telemetry messages evenly distributed since last time telemetry was sent
func (s *deviceSimulator) getNextTelemetryBatch(device *device) *telemetryBatch {
	now := time.Now()
	batchSize := s.simulation.TelemetryBatchSize
	var batch telemetryBatch

	// distribute telemetry messages since the last time telemetry was sent
	// e.g. if telemetry batches of 5 messages are sent at 10:00:00 AM and 10:00:30 AM
	// at 10:00:30 - 5 messages should have creation time of 10:00:10, 10:00:15, 10:00:20, 10:00:25, 10:00:30
	var multiplier int = 0
	interval := s.simulation.TelemetryInterval
	if batchSize > 1 {
		multiplier = ((interval - 1) * 1000) / batchSize
	}

	for i := 0; i < batchSize; i++ {
		creationTime := now.Add(time.Millisecond * time.Duration(-(i * multiplier))) // distribute the messages in the batch evenly
		messages := device.dataGenerator.GenerateTelemetryMessage(device, creationTime)
		batch.messages = append(batch.messages, messages...)
	}

	return &batch
}

// getErrorType groups the error messages into meaningful errors for metrics
func (s *deviceSimulator) getErrorType(err error) string {
	errMsg := "error"
	if err != nil {
		errMsg = err.Error()
	}
	if strings.Contains(errMsg, "429") {
		return "throttled"
	} else if strings.Contains(errMsg, "use of closed connection") || strings.Contains(errMsg, "use of closed network connection") || strings.Contains(errMsg, "forcibly closed by the remote host") {
		return "connection closed"
	} else if strings.Contains(errMsg, "Not Authorized") {
		return "not authorized"
	} else if strings.Contains(errMsg, "context deadline exceeded") {
		return "timeout"
	} else if strings.Contains(reflect.TypeOf(err).String(), "tls.permanentError") {
		return "network error"
	}

	//return errMsg
	return reflect.TypeOf(err).String()
}

func getHubName(connectionString string) string {
	pairs := strings.Split(connectionString, ";")
	for _, pair := range pairs {
		tokens := strings.Split(pair, "=")
		if strings.ToLower(tokens[0]) == "hostname" {
			idx := strings.Index(tokens[1], ".azure-devices.net")
			if idx > 0 {
				return tokens[1][:idx]
			}
		}
	}

	return "unknown"
}
