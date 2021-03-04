package simulating

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"runtime"
	"time"

	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
)

type (
	// Simulator simulates a number of devices connecting to IoT Central
	Simulator struct {
		// cancel function to invoke when the simulation is being stopped.
		cancel context.CancelFunc
		// the context of the simulator
		context context.Context
		// simulator configuration
		config *Config
		// the simulation to execute
		simulation *models.Simulation
		// the target to execute the simulation at
		target *models.SimulationTarget
		// the device configurations to simulate.
		deviceConfigs []*models.SimulationDeviceConfig
		// the models used by the deviceSimulator to simulate.
		models map[string]*models.DeviceModel
		// the devices divides into groups used by the deviceSimulator to simulate.
		deviceGroups map[int]*deviceCollection
		// the device provisioner handling provisioning deviceSimulator.
		provisioner *DeviceProvisioner
		// the device simulator handling simulation of deviceSimulator.
		deviceSimulator *deviceSimulator
	}
)

// Start starts the simulation.
func Start(
	ctx context.Context,
	config *Config,
	simulation *models.Simulation) (*Simulator, error) {

	target, err := storing.Targets.Get(simulation.TargetID)
	if err != nil {
		return nil, err
	}

	deviceConfigs, err := storing.DeviceConfigs.List(simulation.ID)
	if err != nil {
		return nil, err
	}

	deviceModels := map[string]*models.DeviceModel{}
	for _, deviceConfig := range deviceConfigs {
		model, err := storing.DeviceModels.Get(deviceConfig.ModelID)
		if err != nil {
			return nil, err
		}

		deviceModels[model.ID] = model

		simulatedDeviceGauge.WithLabelValues(simulation.ID, simulation.TargetID, deviceConfig.ModelID).Set(float64(deviceConfig.DeviceCount))
	}

	simContext, cancel := context.WithCancel(ctx)
	simulator := &Simulator{
		cancel:          cancel,
		context:         simContext,
		config:          config,
		simulation:      simulation,
		target:          target,
		deviceConfigs:   deviceConfigs,
		models:          deviceModels,
		deviceGroups:    make(map[int]*deviceCollection),
		provisioner:     NewProvisioner(simContext, config),
		deviceSimulator: newDeviceSimulator(simContext, config, simulation),
	}

	// distribute all the devices into groups
	simulator.distributeDeviceGroups()

	// start the simulation pumps
	simulator.start()

	return simulator, nil
}

// start starts the simulation pumps
func (s *Simulator) start() {
	// update the status of simulation
	if err := updateSimulationStatus(s.simulation, models.SimulationStatusStarting); err != nil {
		log.Error().Err(err).Msg("error updating simulation status")
	}

	totalDevices := 0
	for _, dc := range s.deviceConfigs {
		totalDevices += dc.DeviceCount
	}

	log.Debug().
		Bool("enableTelemetry", s.config.EnableTelemetry).
		Bool("enableReportedProps", s.config.EnableReportedProps).
		Bool("enableTwinUpdateAcks", s.config.EnableTwinUpdateAcks).
		Bool("enableCommandAcks", s.config.EnableCommandAcks).
		Int("totalDevices", totalDevices).
		Str("simID", s.simulation.ID).
		Msg("starting simulation")

	// start device simulator
	s.deviceSimulator.start(totalDevices)

	// start telemetry request generator pump
	if s.config.EnableTelemetry {
		go s.startTelemetryRequestPump()
	}

	// start reported props request generator pump
	if s.config.EnableReportedProps {
		go s.startReportedPropertyRequestPump()
	}

	// update the status of simulation
	if err := updateSimulationStatus(s.simulation, models.SimulationStatusRunning); err != nil {
		log.Error().Err(err).Msg("error updating simulation status")
	}
}

// startTelemetryRequestPump starts the pump that sends telemetry requests to devices in waves
func (s *Simulator) startTelemetryRequestPump() {
	log.Debug().Msg("telemetry request generator pump starting")

	for {
		select {
		case <-s.context.Done():
			return
		default:
			// generate a wave of telemetry messages across all device groups
			for waveGroup, devs := range s.deviceGroups {
				select {
				case <-s.context.Done():
					return
				default:
					log.Trace().
						Int("waveGroup", waveGroup).
						Int("numDevices", len(devs.devices)).
						Int("numGoroutines", runtime.NumGoroutine()).
						Msg("sending telemetry requests")

					// send telemetry for all devices in the group
					for _, dev := range devs.devices {
						select {
						case <-s.context.Done():
							return
						default:
							s.deviceSimulator.telemetryRequests <- &telemetryRequest{
								device:  dev,
								context: nil,
							}
						}
					}
					log.Debug().
						Int("waveGroup", waveGroup).
						Int("numDevices", len(devs.devices)).
						Int("numGoroutines", runtime.NumGoroutine()).
						Msg("sent telemetry requests")

					// sleep for some time between each device group
					if waveGroup < s.simulation.WaveGroupCount-1 {
						select {
						case <-s.context.Done():
							return
						case <-time.After(time.Second * time.Duration(s.simulation.WaveGroupInterval)):
						}
					}
				}
			}

			// sleep for some time between each telemetry wave
			select {
			case <-s.context.Done():
				return
			case <-time.After(time.Second * time.Duration(s.simulation.TelemetryInterval)):
			}
		}
	}
}

// startReportedPropertyRequestPump starts the pump that sends reported properties requests to devices in waves
func (s *Simulator) startReportedPropertyRequestPump() {
	// start the reported property update pump after a minute or two to let the device connections settle down
	select {
	case <-s.context.Done():
		return
	case <-time.After(time.Minute * time.Duration(5)):
	}

	log.Debug().Msg("reported properties request generator pump starting")

	for {
		select {
		case <-s.context.Done():
			return
		default:
			// generate a wave of reported property messages across all device groups
			for waveGroup, devs := range s.deviceGroups {
				select {
				case <-s.context.Done():
					return
				default:
					// send reported properties for all devices in the group
					log.Trace().
						Int("waveGroup", waveGroup).
						Int("numDevices", len(devs.devices)).
						Int("numGoroutines", runtime.NumGoroutine()).
						Msg("sending reported properties requests")

					for _, dev := range devs.devices {
						select {
						case <-s.context.Done():
							return
						default:
							s.deviceSimulator.reportedPropsRequests <- &reportedPropsRequest{
								device:  dev,
								context: nil,
							}
						}
					}
					log.Debug().
						Int("waveGroup", waveGroup).
						Int("numDevices", len(devs.devices)).
						Int("numGoroutines", runtime.NumGoroutine()).
						Msg("sent reported properties requests")

					// sleep for some time between each device group
					if waveGroup < s.simulation.WaveGroupCount-1 {
						select {
						case <-s.context.Done():
							return
						case <-time.After(time.Second * time.Duration(s.simulation.WaveGroupInterval)):
						}
					}
				}
			}

			// sleep for some time between each reported property wave
			select {
			case <-s.context.Done():
				return
			case <-time.After(time.Second * time.Duration(s.simulation.ReportedPropsInterval)):
			}
		}
	}
}

// Stop stops the simulation
func (s *Simulator) Stop() error {
	// update the status of simulation
	log.Trace().
		Str("simID", s.simulation.ID).
		Msg("stopping simulation")
	if err := updateSimulationStatus(s.simulation, models.SimulationStatusStopping); err != nil {
		return err
	}

	s.cancel() // send cancellation signal to all go funcs.

	s.deviceSimulator.stop()

	// disconnect all devices
	for _, devs := range s.deviceGroups {
		for _, dev := range devs.devices {
			s.deviceSimulator.disconnectDevice(dev)
		}
	}

	// update the status of simulation
	if err := updateSimulationStatus(s.simulation, models.SimulationStatusStopped); err != nil {
		return err
	}

	log.Debug().
		Str("simID", s.simulation.ID).
		Msg("simulation stopped")

	return nil
}

// distributeDeviceGroups divides the devices in the simulation into wave groups
func (s *Simulator) distributeDeviceGroups() {
	var totalDevices int = 0
	for _, dc := range s.deviceConfigs {
		totalDevices += dc.DeviceCount
	}

	// divide total devices in simulation into wave groups
	waveGroupCount := s.simulation.WaveGroupCount
	devicesAdded := 0
	devicesPerWave := totalDevices
	if waveGroupCount > 0 && waveGroupCount < totalDevices {
		devicesPerWave = totalDevices / waveGroupCount
	}

	// go over all device models and divide all devices into wave groups based on above calculations
	for _, deviceCfg := range s.deviceConfigs {
		model := s.models[deviceCfg.ModelID]
		for i := 1; i <= deviceCfg.DeviceCount; i++ {
			deviceID := fmt.Sprintf("%s-%s-%s-%d",
				s.simulation.ID,
				s.target.ID,
				deviceCfg.ID,
				i)

			group := devicesAdded / devicesPerWave
			if group > waveGroupCount {
				// handle odd cases; e.g.: total devices = 7, wave groups = 2, device #7 should be part of wave group 2
				group--
			}

			if _, found := s.deviceGroups[group]; found == false {
				s.deviceGroups[group] = new(deviceCollection)
			}

			deviceContext, deviceCancel := context.WithCancel(s.context)
			d := device{
				deviceID:             deviceID,
				model:                model,
				target:               s.target,
				connectionString:     "",
				isConnected:          false,
				isConnecting:         false,
				telemetrySentTime:    time.Time{},
				sendingTelemetry:     false,
				sendingReportedProps: false,
				iotHubClient:         nil,
				twinSub:              nil,
				c2dSub:               nil,
				retryCount:           0,
				cancel:               deviceCancel,
				context:              deviceContext,
				dataGenerator: &DataGenerator{
					CapabilityModel: model.ParseDeviceCapabilityModel(),
				},
			}
			s.deviceGroups[group].devices = append(s.deviceGroups[group].devices, &d)
			devicesAdded++
		}
	}
}

// sleep sleeps for the given duration with cancellation context
func sleep(ctx context.Context, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	case <-time.After(duration):
	}
}

// randSleep sleeps for random time within the given min/max range with cancellation context
func randSleep(ctx context.Context, minMs int, maxMs int) {
	sleep(ctx, time.Millisecond*time.Duration(minMs+rand.Intn(maxMs)))
}

func updateSimulationStatus(simulation *models.Simulation, status models.SimulationStatus) error {
	// update the status of simulation
	simulation.Status = status
	err := storing.Simulations.Set(simulation)
	return err
}
