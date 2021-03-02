package serving

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
)

// listDeviceConfigs lists all device configurations for a simulation.
func listDeviceConfigs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]

	items, err := storing.DeviceConfigs.List(simID)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	handleError(err, w)
}

// getDeviceConfig gets a specific device configuration for a simulation.
func getDeviceConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]
	cfgID := vars["configId"]

	item, err := storing.DeviceConfigs.Get(simID, cfgID)
	if handleError(err, w) {
		return
	}

	if item == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(item)
	handleError(err, w)
}

// upsertDeviceConfig creates a new or updates an existing device configuration for a simulation.
func upsertDeviceConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]

	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var cfg models.SimulationDeviceConfig
	err = json.Unmarshal(req, &cfg)
	if handleError(err, w) {
		return
	}

	err = storing.DeviceConfigs.Set(simID, &cfg)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cfg)
	handleError(err, w)
}

// deleteDeviceConfig deletes an existing device configuration for a simulation.
func deleteDeviceConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	simID := vars["id"]
	cfgID := vars["configId"]
	err := storing.DeviceConfigs.Delete(simID, cfgID)
	handleError(err, w)
}
