package controlling

import (
	"context"
	"fmt"
	"github.com/iot-for-all/starling/pkg/config"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/simulating"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/rs/zerolog/log"
	"net/http"
	"sync"
	"time"
)

// Controller responsible for starting and stopping simulations; provisioning and deleting devices from a target application.
type Controller struct {
	context     context.Context      // parent program context.
	globalCfg   *config.GlobalConfig // global configuration.
	simulations map[string]*simulating.Simulator
}

// NewController creates a new controller.
func NewController(context context.Context, globalConfig *config.GlobalConfig) *Controller {
	return &Controller{
		context:     context,
		globalCfg:   globalConfig,
		simulations: map[string]*simulating.Simulator{},
	}
}

// StartSimulation starts a simulation.
func (c *Controller) StartSimulation(simulation *models.Simulation) error {
	if _, ok := c.simulations[simulation.ID]; ok == true {
		return fmt.Errorf("simulation %s is already running. stop it first and then try running it again", simulation.ID)
	}

	simulator, err := simulating.Start(c.context, &c.globalCfg.Simulation, simulation)
	if err != nil {
		return err
	}

	c.simulations[simulation.ID] = simulator
	return nil
}

// StopSimulation stops a simulation.
func (c *Controller) StopSimulation(simulation *models.Simulation) error {
	sim, ok := c.simulations[simulation.ID]
	if !ok {
		return fmt.Errorf("simulation %s is not running. nothing to stop", simulation.ID)
	}

	if err := sim.Stop(); err != nil {
		return err
	}

	delete(c.simulations, simulation.ID)
	return nil
}

// ProvisionDevices provisions devices in a target based on the deviceConfig.
func (c *Controller) ProvisionDevices(ctx context.Context, simulation *models.Simulation, target *models.SimulationTarget, model *models.DeviceModel, maxDeviceID int, numDevices int) error {
	provisioner := simulating.NewProvisioner(c.context, &c.globalCfg.Simulation)

	wg := sync.WaitGroup{}
	for i := 1; i <= numDevices; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			deviceID := fmt.Sprintf("%s-%s-%s-%d",
				simulation.ID,
				target.ID,
				model.ID,
				maxDeviceID+i)

			wg.Add(1)
			go c.provisionDevice(simulation, target, model, deviceID, provisioner, &wg)

			// throttle DPS registrations
			if i%c.globalCfg.Simulation.MaxConcurrentRegistrations == 0 {
				wg.Wait()
			}

			if i%10 == 0 {
				log.Debug().
					Int("provisioned", i).
					Int("remaining", numDevices-i).
					Str("modelID", model.ID).
					Msg("provisioning in progress")
			}
		}
	}
	wg.Wait()

	log.Debug().
		Int("provisioned", numDevices).
		Str("modelID", model.ID).
		Msg("provisioning completed")

	return nil
}

// provisionDevice provisions a device in IoT Central and saves it into the database cache.
func (c *Controller) provisionDevice(simulation *models.Simulation, target *models.SimulationTarget,
	model *models.DeviceModel, deviceID string, provisioner *simulating.DeviceProvisioner,
	wg *sync.WaitGroup) {
	defer wg.Done()

	req := &simulating.ProvisioningRequest{
		DeviceID:   deviceID,
		Context:    c.context,
		Target:     target,
		Simulation: simulation,
		Model:      model,
	}

	// call DPS to register the device
	result := provisioner.Provision(req)
	if result == nil {
		return
	}

	// cache the device for future use
	newDevice := models.SimulationTargetDevice{
		TargetID:         req.Target.ID,
		DeviceID:         req.DeviceID,
		ConnectionString: result.ConnectionString,
	}
	storing.TargetDevices.Set(&newDevice)
	//log.Debug().Str("deviceId", req.DeviceID).Msg("saved device in cache")
}

func (c *Controller) DeleteAllDevices(ctx context.Context, sim *models.Simulation, target *models.SimulationTarget) error {
	targetDevices, err := storing.TargetDevices.ListByTargetIdSimId(target.ID, sim.ID)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	numDevices := len(targetDevices)
	for i, td := range targetDevices {
		select {
		case <-ctx.Done():
			return nil
		default:

			wg.Add(1)
			go c.deleteDevice(ctx, target, td.DeviceID, &wg)

			// throttle API calls to Central
			if i%c.globalCfg.Simulation.MaxConcurrentDeletes == 0 {
				wg.Wait()
			}

			if i%10 == 0 {
				log.Debug().
					Int("deleted", i+1).
					Int("remaining", numDevices-i-1).
					Str("simId", sim.ID).
					Msg("deleting devices in progress")
			}
		}
	}
	wg.Wait()

	log.Debug().
		Int("deleted", numDevices).
		Str("simId", sim.ID).
		Msg("deleting devices completed")

	return nil

}

// DeleteDevices deletes devices in a target based on the deviceConfig.
func (c *Controller) DeleteDevices(ctx context.Context, simulation *models.Simulation, target *models.SimulationTarget, model *models.DeviceModel, maxDeviceID int, numDevices int) error {
	wg := sync.WaitGroup{}

	for i := 0; i < numDevices; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			deviceID := fmt.Sprintf("%s-%s-%s-%d",
				simulation.ID,
				target.ID,
				model.ID,
				maxDeviceID-i)

			wg.Add(1)
			go c.deleteDevice(ctx, target, deviceID, &wg)

			// throttle API calls to Central
			if i%c.globalCfg.Simulation.MaxConcurrentDeletes == 0 {
				wg.Wait()
			}

			if i%10 == 0 {
				log.Debug().
					Int("deleted", i+1).
					Int("remaining", numDevices-i-1).
					Str("modelID", model.ID).
					Msg("deleting devices in progress")
			}
		}
	}
	wg.Wait()

	log.Debug().
		Int("deleted", numDevices).
		Str("modelID", model.ID).
		Msg("deleting devices completed")

	return nil
}

// deleteDevice deletes a device from the target application and local database cache
func (c *Controller) deleteDevice(ctx context.Context, target *models.SimulationTarget, deviceID string, wg *sync.WaitGroup) {
	defer wg.Done()

	path := fmt.Sprintf("https://%s/api/devices/%s?api-version=1.0", target.AppUrl, deviceID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", path, nil)
	if err != nil {
		log.Err(err).Str("deviceID", deviceID).Msg("error creating deprovision request")
		return
	}

	req.Header.Add("Authorization", target.AppToken)
	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		log.Err(err).Str("deviceID", deviceID).Msg("error executing delete request")
		// ignore delete errors
	}

	// user might want to delete devices from Central that might not exist in client side case
	// ignore errors
	_ = storing.TargetDevices.Delete(target.ID, deviceID)

	log.Trace().Str("deviceID", deviceID).Str("path", path).Msg("deleted device")
}

// ResetSimulationStatus resets all simulation status to stopped
func (c *Controller) ResetSimulationStatus() error {
	sims, err := storing.Simulations.List()
	if err != nil {
		return err
	}

	for _, sim := range sims {
		sim.Status = models.SimulationStatusReady
		sim.LastUpdatedTime = time.Now()
		err := storing.Simulations.Set(&sim)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetConnectedDeviceCount returns the number of devices connected for the given simulation and model
func (c *Controller) GetConnectedDeviceCount(simulation *models.Simulation, modelId string) int {
	sim, ok := c.simulations[simulation.ID]
	if !ok {
		return 0
	}

	return sim.GetConnectedDeviceCount(modelId)
}

func (c *Controller) GetMetricsStatus(ctx context.Context) models.MetricsStatus {
	return models.MetricsStatus{
		GrafanaServer:    c.getServerStatus(ctx, "Grafana", c.globalCfg.HTTP.GrafanaPort),
		PrometheusServer: c.getServerStatus(ctx, "Prometheus", c.globalCfg.HTTP.PrometheusPort),
	}
}

func (c *Controller) getServerStatus(ctx context.Context, name string, portNumber int) bool {
	path := fmt.Sprintf("http://localhost:%d", portNumber)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		log.Err(err).Str("path", path).Msg(fmt.Sprintf("error creating %s server status request", name))
		return false
	}

	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	_, err = client.Do(req)
	if err != nil {
		log.Err(err).Str("path", path).Msg("error getting Grafana server status request")
		return false
	}
	return true
}
