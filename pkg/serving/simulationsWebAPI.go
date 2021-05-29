package serving

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// webAPIListSimulations lists all simulations.
func webAPIListSimulations(w http.ResponseWriter, r *http.Request) {
	sims, err := storing.Simulations.List()
	if handleError(err, w) {
		return
	}

	deviceModels, err := storing.DeviceModels.List()
	if handleError(err, w) {
		return
	}

	simViews := make([]models.SimulationView, len(sims))
	for index, sim := range sims {
		configs, err := storing.DeviceConfigs.List(sim.ID)
		if handleError(err, w) {
			return
		}

		deviceViews := make([]models.SimulationViewDeviceConfig, len(configs))
		for i, config := range configs {
			deviceViews[i].ID = config.ID
			deviceViews[i].ModelID = config.ModelID
			deviceViews[i].SimulatedCount = config.DeviceCount
			provisionedCount, err := getProvisionedDeviceCount(sim.ID, sim.TargetID, config.ModelID)
			if handleError(err, w) {
				return
			}
			deviceViews[i].ProvisionedCount = provisionedCount
			// get the connected device count from prometheus
			deviceViews[i].ConnectedCount = controller.GetConnectedDeviceCount(&sim, config.ModelID)
		}

		// add device models which might have been added since this simulation is created
		for _, dm := range deviceModels {
			foundDM := false
			for _, dv := range deviceViews {
				if dv.ModelID == dm.ID {
					foundDM = true
					break
				}
			}

			if !foundDM {
				deviceViews = append(deviceViews, models.SimulationViewDeviceConfig{
					ID:               dm.ID,
					ModelID:          dm.ID,
					ProvisionedCount: 0,
					SimulatedCount:   0,
					ConnectedCount:   0,
				})
			}
		}

		simViews[index] = models.SimulationView{
			Simulation: sim,
			Devices:    deviceViews,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(simViews)
	handleError(err, w)
}

// webAPIGetSimulation gets an existing simulation by its id.
func webAPIGetSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sim, err := storing.Simulations.Get(id)
	if handleError(err, w) {
		return
	}

	if sim == nil {
		http.NotFound(w, r)
		return
	}

	deviceModels, err := storing.DeviceModels.List()
	if handleError(err, w) {
		return
	}

	configs, err := storing.DeviceConfigs.List(id)
	if handleError(err, w) {
		return
	}

	deviceViews := make([]models.SimulationViewDeviceConfig, len(configs))
	for i, config := range configs {
		deviceViews[i].ID = config.ID
		deviceViews[i].ModelID = config.ModelID
		deviceViews[i].SimulatedCount = config.DeviceCount
		provisionedCount, err := getProvisionedDeviceCount(sim.ID, sim.TargetID, config.ModelID)
		if handleError(err, w) {
			return
		}
		deviceViews[i].ProvisionedCount = provisionedCount
		// get the connected device count from prometheus
		deviceViews[i].ConnectedCount = controller.GetConnectedDeviceCount(sim, config.ModelID)
	}

	// add device models which might have been added since this simulation is created
	for _, dm := range deviceModels {
		foundDM := false
		for _, dv := range deviceViews {
			if dv.ModelID == dm.ID {
				foundDM = true
				break
			}
		}

		if !foundDM {
			deviceViews = append(deviceViews, models.SimulationViewDeviceConfig{
				ID:               dm.ID,
				ModelID:          dm.ID,
				ProvisionedCount: 0,
				SimulatedCount:   0,
				ConnectedCount:   0,
			})
		}
	}

	simView := models.SimulationView{
		Simulation: *sim,
		Devices:    deviceViews,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(simView)
	handleError(err, w)
}

func getProvisionedDeviceCount(simID string, targetID string, modelID string) (int, error) {
	targetDevices, err := storing.TargetDevices.List(targetID)
	if err != nil {
		return 0, err
	}

	prefix := fmt.Sprintf("%s-%s-%s-", simID, targetID, modelID)
	numProvisioned := 0
	for _, d := range targetDevices {
		if strings.Index(d.DeviceID, prefix) == 0 {
			numProvisioned++
		}
	}

	return numProvisioned, nil
}

// webAPIAddSimulation add a new simulation.
func webAPIAddSimulation(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var simView models.SimulationView
	err = json.Unmarshal(req, &simView)
	if handleError(err, w) {
		return
	}

	// check to see if there is an existing sim with same ID
	existingSim, err := storing.Simulations.Get(simView.ID)
	if handleError(err, w) {
		return
	}

	if existingSim != nil {
		msg := fmt.Sprintf("Another simulation exists with Simulation ID '%s'. Try again with another ID.", simView.ID)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// reset status
	simView.Simulation.Status = models.SimulationStatusReady
	simView.Simulation.LastUpdatedTime = time.Now()

	// save simulation
	err = storing.Simulations.Set(&simView.Simulation)
	if handleError(err, w) {
		return
	}

	// save new simulation device configs
	for _, simViewDeviceConfig := range simView.Devices {
		dc := &models.SimulationDeviceConfig{
			ID:          simViewDeviceConfig.ID,
			ModelID:     simViewDeviceConfig.ModelID,
			DeviceCount: simViewDeviceConfig.SimulatedCount,
		}
		err = storing.DeviceConfigs.Set(simView.ID, dc)
		if handleError(err, w) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&simView)
	handleError(err, w)
}

// webAPIUpdateSimulation update an existing simulation.
func webAPIUpdateSimulation(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var simView models.SimulationView
	err = json.Unmarshal(req, &simView)
	if handleError(err, w) {
		return
	}

	// check simulation status
	sim, err := storing.Simulations.Get(simView.ID)
	if handleError(err, w) {
		return
	}

	if sim != nil && sim.Status != models.SimulationStatusReady {
		msg := fmt.Sprintf("Simulation cannot be updated while it is in '%s' status .", sim.Status)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// update lastUpdatedDate
	simView.Simulation.LastUpdatedTime = time.Now()

	// update simulation
	err = storing.Simulations.Set(&simView.Simulation)
	if handleError(err, w) {
		return
	}

	// delete old simulation device configs
	deviceConfigs, err := storing.DeviceConfigs.List(simView.ID)
	if handleError(err, w) {
		return
	}
	for _, dc := range deviceConfigs {
		err := storing.DeviceConfigs.Delete(simView.ID, dc.ID)
		if handleError(err, w) {
			return
		}
	}

	// add new simulation device configs
	for _, simViewDeviceConfig := range simView.Devices {
		dc := &models.SimulationDeviceConfig{
			ID:          simViewDeviceConfig.ID,
			ModelID:     simViewDeviceConfig.ModelID,
			DeviceCount: simViewDeviceConfig.SimulatedCount,
		}
		err = storing.DeviceConfigs.Set(simView.ID, dc)
		if handleError(err, w) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&simView)
	handleError(err, w)
}

// webAPIDeleteSimulation deletes an existing simulation.
func webAPIDeleteSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sim, err := storing.Simulations.Get(id)
	if handleError(err, w) {
		return
	}

	// Simulation cannot be deleted when running.
	if sim.Status != models.SimulationStatusReady {
		msg := fmt.Sprintf("Simulation cannot be deleted while it is in '%s' status.", sim.Status)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	target, err := storing.Targets.Get(sim.TargetID)
	if handleError(err, w) {
		return
	}

	sim.Status = models.SimulationStatusDeleting
	if err := storing.Simulations.Set(sim); err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to deleting status")
		return
	}

	log.Debug().Msg("starting parallel delete")
	// de-provision and delete all target devices for this simulation in background
	go deleteSimulationInternal(sim, target)
	log.Debug().Msg("returning response")

	w.Write([]byte("simulation is getting deleted in background"))
}

// webAPIStartSimulation starts an existing simulation.
func webAPIStartSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sim, err := storing.Simulations.Get(id)
	if handleError(err, w) {
		return
	}

	if sim == nil {
		http.NotFound(w, r)
		return
	}

	if sim.Status != models.SimulationStatusReady {
		msg := fmt.Sprintf("Simulation cannot be started while it is in '%s' status.", sim.Status)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = controller.StartSimulation(sim)
	handleError(err, w)
}

// webAPIStopSimulation stops a running simulation.
func webAPIStopSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sim, err := storing.Simulations.Get(id)
	if handleError(err, w) {
		return
	}

	if sim.Status != models.SimulationStatusRunning {
		msg := fmt.Sprintf("Simulation cannot be stopped while it is in '%s' status.", sim.Status)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = controller.StopSimulation(sim)
	handleError(err, w)
}

// webAPIProvisionDevices provisions devices in a target based on the device configs from the given start index
func webAPIProvisionDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]

	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var simViewDeviceConfigs []models.SimulationViewDeviceConfig
	err = json.Unmarshal(req, &simViewDeviceConfigs)
	if handleError(err, w) {
		return
	}

	sim, err := storing.Simulations.Get(simID)
	if handleError(err, w) {
		return
	}

	if sim.Status != models.SimulationStatusReady {
		msg := fmt.Sprintf("Devices cannot be provisioned while the simulation is in '%s' status.", sim.Status)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	target, err := storing.Targets.Get(sim.TargetID)
	if handleError(err, w) {
		return
	}

	// kick off device provisioning in the background
	go provisionDevicesInternal(sim, target, simViewDeviceConfigs)

	w.Write([]byte("device provisioning started in background"))
}

// webAPIExportSimulation export the simulation as a shell script
func webAPIExportSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]

	sim, err := storing.Simulations.Get(simID)
	if handleError(err, w) {
		return
	}
	if sim == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/x-sh")
	w.Header().Set("Content-Disposition", "filename=loadData-"+sim.ID+".sh")

	sim.Status = models.SimulationStatusReady
	exporter := exporter{sim: sim}
	exportFileContent, err := exporter.exportSimulation()
	if handleError(err, w) {
		return
	}

	w.Write(exportFileContent)
	handleError(err, w)
}

func deleteSimulationInternal(sim *models.Simulation, target *models.SimulationTarget) {
	ctx := context.Background()
	if err := controller.DeleteAllDevices(ctx, sim, target); err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error deleting all devices")

		sim.Status = models.SimulationStatusReady
		if err := storing.Simulations.Set(sim); err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
			return
		}
		return
	}

	// delete all device config
	deviceConfigs, err := storing.DeviceConfigs.List(sim.ID)
	if err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error deleting all devices that belong to simulation")
		sim.Status = models.SimulationStatusReady
		if err := storing.Simulations.Set(sim); err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
			return
		}
		return
	}

	for _, dc := range deviceConfigs {
		err = storing.DeviceConfigs.Delete(sim.ID, dc.ID)
		if err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Msg("error deleting device configs that belong to simulation")
			sim.Status = models.SimulationStatusReady
			if err := storing.Simulations.Set(sim); err != nil {
				log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
				return
			}
			return
		}
	}

	// delete simulation
	err = storing.Simulations.Delete(sim.ID)
	if err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error deleting simulation")
		sim.Status = models.SimulationStatusReady
		if err := storing.Simulations.Set(sim); err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
			return
		}
		return
	}
	log.Debug().Str("simID", sim.ID).Msg("deleted simulation")
}

func provisionDevicesInternal(sim *models.Simulation, target *models.SimulationTarget, simViewDeviceConfigs []models.SimulationViewDeviceConfig) {
	ctx := context.Background()

	sim.Status = models.SimulationStatusProvisioning
	if err := storing.Simulations.Set(sim); err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to provisioning status")
		return
	}

	for _, svdc := range simViewDeviceConfigs {

		model, err := storing.DeviceModels.Get(svdc.ModelID)
		if err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Str("model", svdc.ModelID).Msg("error provisioning devices, could not find model")

			sim.Status = models.SimulationStatusReady
			if err := storing.Simulations.Set(sim); err != nil {
				log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
			}
			return
		}

		targetDevices, err := storing.TargetDevices.List(target.ID)
		if err != nil {
			log.Error().Err(err).Str("simID", sim.ID).Str("target", target.ID).Msg("error provisioning devices, could not find target devices")
			sim.Status = models.SimulationStatusReady
			if err := storing.Simulations.Set(sim); err != nil {
				log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
			}
			return
		}

		numDevicesRequested := svdc.ProvisionedCount

		// get the maximum device ID for this model
		maxDeviceID := 0
		numDevicesExisting := 0
		if targetDevices != nil {
			// format SimID-TargetID-modelID-NNNN
			prefix := fmt.Sprintf("%s-%s-%s-", sim.ID, target.ID, svdc.ModelID)
			length := len(prefix)
			for _, d := range targetDevices {
				if strings.Index(d.DeviceID, prefix) == 0 {
					idStr := d.DeviceID[length:]
					did, _ := strconv.Atoi(idStr)
					if maxDeviceID < did {
						maxDeviceID = did
					}
					numDevicesExisting++
				}
			}
		}

		if numDevicesRequested > numDevicesExisting {
			// overall more devices are requested, so provision the difference
			numDevicesToProvision := numDevicesRequested - numDevicesExisting
			if err := controller.ProvisionDevices(ctx, sim, target, model, maxDeviceID, numDevicesToProvision); err != nil {
				log.Error().Err(err).Str("simID", sim.ID).Str("target", target.ID).Int("numDevicesExisting", numDevicesExisting).
					Int("numDevicesRequested", numDevicesRequested).Str("model", model.ID).Msg("error provisioning devices")
			}
		} else if numDevicesRequested < numDevicesExisting {
			// overall less devices are requested, so delete the excess devices
			numDevicesToDelete := numDevicesExisting - numDevicesRequested
			if err := controller.DeleteDevices(ctx, sim, target, model, maxDeviceID, numDevicesToDelete); err != nil {
				log.Error().Err(err).Str("simID", sim.ID).Str("target", target.ID).Int("numDevicesExisting", numDevicesExisting).
					Int("numDevicesRequested", numDevicesRequested).Str("model", model.ID).Msg("error deleting devices")
			}
		} else {
			// old and new counts are same, ignore the request for this model
			log.Debug().Str("simID", sim.ID).Str("target", target.ID).
				Str("model", model.ID).Int("numDevicesExisting", numDevicesExisting).
				Int("numDevicesRequested", numDevicesRequested).Msg("ignoring provision devices as there is nothing to do")
		}
	}

	sim.Status = models.SimulationStatusReady
	if err := storing.Simulations.Set(sim); err != nil {
		log.Error().Err(err).Str("simID", sim.ID).Msg("error updating simulation to ready status")
		return
	}
}
