package serving

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"io/ioutil"
	"net/http"
)

// listTargets lists all targets.
func listTargets(w http.ResponseWriter, _ *http.Request) {
	items, err := storing.Targets.List()
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	handleError(err, w)
}

// listTargetDevicesByTarget lists all devices in a target.
func listTargetDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	items, err := storing.TargetDevices.List(id)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	handleError(err, w)
}

// getTarget gets a target by its id.
func getTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	t, err := storing.Targets.Get(id)
	if handleError(err, w) {
		return
	}

	if t == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(t)
	handleError(err, w)
}

// listTargetModels lists all models for a target.
func getTargetModels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := storing.TargetModels.Get(id)
	if handleError(err, w) {
		return
	}

	if m == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(m.Models)
	handleError(err, w)
}

// getTargetDevice gets a device in a target by its target id and deviceId
func getTargetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deviceId := vars["deviceId"]
	d, err := storing.TargetDevices.Get(id, deviceId)
	if handleError(err, w) {
		return
	}

	if d == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(d)
	handleError(err, w)
}

// upsertTarget adds a new or updates an existing target
func upsertTarget(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var t models.SimulationTarget
	err = json.Unmarshal(req, &t)
	if handleError(err, w) {
		return
	}

	upsertTargetInternal(w, r, t)
}

func upsertTargetInternal(w http.ResponseWriter, r *http.Request, t models.SimulationTarget) {
	err := storing.Targets.Set(&t)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&t)
	handleError(err, w)
}

// upsertTargetModels adds or updates existing models for a target.
func upsertTargetModels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var m []string
	err = json.Unmarshal(req, &m)
	if handleError(err, w) {
		return
	}

	tm := models.SimulationTargetModels{
		TargetID: id,
		Models:   m,
	}

	err = storing.TargetModels.Set(&tm)
	handleError(err, w)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&tm)
	handleError(err, w)
}

// upsertTargetDevice adds or updates existing device in a target.
func upsertTargetDevice(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var d models.SimulationTargetDevice
	err = json.Unmarshal(req, &d)
	if handleError(err, w) {
		return
	}

	err = storing.TargetDevices.Set(&d)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&d)
	handleError(err, w)
}

// deleteTarget deletes an existing target
func deleteTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storing.Targets.Delete(id)
	if handleError(err, w) {
		return
	}

	err = storing.TargetModels.Delete(id)
	handleError(err, w)
}

// deleteTargetModels deletes existing target models
func deleteTargetModels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storing.TargetModels.Delete(id)
	handleError(err, w)
}

// deleteTargetDevice deletes existing target device
func deleteTargetDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deviceId := vars["deviceId"]
	err := storing.TargetDevices.Delete(id, deviceId)
	handleError(err, w)
}

// deleteAllTargetDevices deletes all devices from the target
func deleteAllTargetDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storing.TargetDevices.DeleteAll(id)
	handleError(err, w)
}
