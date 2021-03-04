package controlling

import (
	"context"
	"fmt"
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
	context       context.Context    // parent program context.
	simulationCfg *simulating.Config // simulator configuration.
	simulations   map[string]*simulating.Simulator
}

// NewController creates a new controller.
func NewController(context context.Context, simulatorCfg *simulating.Config) *Controller {
	return &Controller{
		context:       context,
		simulationCfg: simulatorCfg,
		simulations:   map[string]*simulating.Simulator{},
	}
}

// StartSimulation starts a simulation.
func (c *Controller) StartSimulation(simulation *models.Simulation) error {
	if _, ok := c.simulations[simulation.ID]; ok == true {
		return fmt.Errorf("simulation %s is already running. stop it first and then try running it again", simulation.ID)
	}

	simulator, err := simulating.Start(c.context, c.simulationCfg, simulation)
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
	provisioner := simulating.NewProvisioner(c.context, c.simulationCfg)

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
			if i%c.simulationCfg.MaxConcurrentRegistrations == 0 {
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
			if i%c.simulationCfg.MaxConcurrentDeletes == 0 {
				wg.Wait()
			}

			if i%10 == 0 {
				log.Debug().
					Int("deleted", i).
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

	path := fmt.Sprintf("https://%s/api/preview/devices/%s", target.AppUrl, deviceID)
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
}

// ResetSimulationStatus resets all simulation status to stopped
func (c *Controller) ResetSimulationStatus() error {
	sims, err := storing.Simulations.List()
	if err != nil {
		return err
	}

	for _, sim := range sims {
		sim.Status = models.SimulationStatusStopped
		err := storing.Simulations.Set(&sim)
		if err != nil {
			return err
		}
	}
	return nil
}
