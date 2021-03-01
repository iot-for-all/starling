package serving

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/reddyduggempudi/starling/pkg/models"
	"github.com/reddyduggempudi/starling/pkg/storing"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// simulationDetail represents simulation along with device configurations.
type simulationDetail struct {
	models.Simulation
	DeviceConfig []*models.SimulationDeviceConfig `json:"deviceConfig"`
}

// listSimulations lists all simulations.
func listSimulations(w http.ResponseWriter, _ *http.Request) {
	items, err := storing.Simulations.List()
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	handleError(err, w)
}

// getSimulation gets an existing simulation by its id.
func getSimulation(w http.ResponseWriter, r *http.Request) {
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

	configs, err := storing.DeviceConfigs.List(id)
	if handleError(err, w) {
		return
	}

	o := simulationDetail{
		Simulation:   *sim,
		DeviceConfig: configs,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(o)
	handleError(err, w)
}

// upsertSimulation add a new simulation or update an existing simulation.
func upsertSimulation(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var sim models.Simulation
	err = json.Unmarshal(req, &sim)
	if handleError(err, w) {
		return
	}

	err = storing.Simulations.Set(&sim)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&sim)
	handleError(err, w)
}

// deleteSimulation deletes an existing simulation.
func deleteSimulation(w http.ResponseWriter, r *http.Request) {
	// TODO: Simulation cannot be deleted when running.
	vars := mux.Vars(r)
	id := vars["id"]
	err := storing.Simulations.Delete(id)
	handleError(err, w)
}

// startSimulation starts an existing simulation.
func startSimulation(w http.ResponseWriter, r *http.Request) {
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

	// todo: handle error
	err = controller.StartSimulation(sim)
	handleError(err, w)
}

// stopSimulation stops a running simulation.
func stopSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sim, err := storing.Simulations.Get(id)
	if handleError(err, w) {
		return
	}

	err = controller.StopSimulation(sim)
	handleError(err, w)
}

// provisionDevices provisions devices in a target based on the device configs from the given start index
func provisionDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]
	modelID := vars["modelId"]
	numDevices, _ := strconv.Atoi(vars["numDevices"])
	if numDevices < 1 {
		return
	}

	sim, err := storing.Simulations.Get(simID)
	if handleError(err, w) {
		return
	}

	target, err := storing.Targets.Get(sim.TargetID)
	if handleError(err, w) {
		return
	}

	model, err := storing.DeviceModels.Get(modelID)
	if handleError(err, w) {
		return
	}

	targetDevices, err := storing.TargetDevices.List(target.ID)
	if handleError(err, w) {
		return
	}

	maxDeviceID := 0
	if targetDevices != nil {
		for _, d := range targetDevices {
			// format SimID-AppID-modelID-NNNN
			tokens := strings.Split(d.DeviceID, "-")
			if modelID == tokens[2] {
				did, _ := strconv.Atoi(tokens[3])
				if maxDeviceID < did {
					maxDeviceID = did
				}
			}
		}
	}

	if err := controller.ProvisionDevices(r.Context(), sim, target, model, maxDeviceID, numDevices); err != nil {
		if handleError(err, w) {
			return
		}
	}
}

// deleteDevices deletes the devices from the target and local cache based on the device configs from the given start index
func deleteDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]
	modelID := vars["modelId"]
	numDevices, _ := strconv.Atoi(vars["numDevices"])
	if numDevices < 1 {
		return
	}

	sim, err := storing.Simulations.Get(simID)
	if handleError(err, w) {
		return
	}

	target, err := storing.Targets.Get(sim.TargetID)
	if handleError(err, w) {
		return
	}

	model, err := storing.DeviceModels.Get(modelID)
	if handleError(err, w) {
		return
	}

	targetDevices, err := storing.TargetDevices.List(target.ID)
	if handleError(err, w) {
		return
	}

	maxDeviceID := 0
	if targetDevices != nil {
		for _, d := range targetDevices {
			// format SimID-AppID-modelID-NNNN
			tokens := strings.Split(d.DeviceID, "-")
			if modelID == tokens[2] {
				did, _ := strconv.Atoi(tokens[3])
				if maxDeviceID < did {
					maxDeviceID = did
				}
			}
		}
	}

	if err := controller.DeleteDevices(r.Context(), sim, target, model, maxDeviceID, numDevices); err != nil {
		if handleError(err, w) {
			return
		}
	}
}
