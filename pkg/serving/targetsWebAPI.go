package serving

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

// webAPIListTargets lists all targets.
func webAPIListTargets(w http.ResponseWriter, r *http.Request) {
	listTargets(w, r)
}

// webAPIGetTarget gets a target by its id.
func webAPIGetTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		msg := fmt.Sprintf("application id is required for GET targets call.")
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	getTarget(w, r)
}

// webAPIAddTarget adds a new or updates an existing target
func webAPIAddTarget(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var t models.SimulationTarget
	var tv models.SimulationTargetView
	err = json.Unmarshal(req, &t)
	if handleError(err, w) {
		return
	}
	err = json.Unmarshal(req, &tv)
	if handleError(err, w) {
		return
	}

	// check to see if there are existing targets with same id
	target, err := storing.Targets.Get(t.ID)
	if handleError(err, w) {
		return
	}

	if target != nil {
		msg := fmt.Sprintf("Another application exists with Application ID '%s'. Try again with another ID.", t.ID)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	upsertTargetInternal(w, r, t)

	// download models and store them into database
	if tv.ImportModels {
		dtDownloader := NewDeviceTemplateDownloader(&t)
		deviceModels, err := dtDownloader.DownloadModels()
		if handleError(err, w) {
			return
		}

		for _, dm := range deviceModels {
			err := storing.DeviceModels.Set(dm)
			if handleError(err, w) {
				return
			}
		}
	}

	// add model bindings with target
	deviceModels, err := storing.DeviceModels.List()
	if handleError(err, w) {
		return
	}

	var tm models.SimulationTargetModels
	tm.TargetID = t.ID
	tm.Models = make([]string, len(deviceModels))
	for i, dm := range deviceModels {
		tm.Models[i] = dm.ID
	}

	err = storing.TargetModels.Set(&tm)
	handleError(err, w)
}

// webAPIUpdateTarget updates updates an existing target
func webAPIUpdateTarget(w http.ResponseWriter, r *http.Request) {
	upsertTarget(w, r)
}

// webAPIDeleteTarget deletes an existing target
func webAPIDeleteTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// check if there are any simulations using this target
	sims, err := storing.Simulations.List()
	if handleError(err, w) {
		return
	}

	for _, sim := range sims {
		if sim.TargetID == id {
			msg := fmt.Sprintf("application '%s' cannot be deleted as there are simulations using this model. Delete the simulation '%s' and try again.", id, sim.Name)
			log.Error().Msg(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	}

	// delete model bindings with target
	err = storing.TargetModels.Delete(id)
	if handleError(err, w) {
		return
	}

	deleteTarget(w, r)
}

// webAPIImportModelsFromTarget imports device models from an existing target
func webAPIImportModelsFromTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	target, err := storing.Targets.Get(id)
	if handleError(err, w) {
		return
	}

	// download device models from the application
	dtDownloader := NewDeviceTemplateDownloader(target)
	deviceModels, err := dtDownloader.DownloadModels()
	if handleError(err, w) {
		return
	}

	// store the models in the database
	for _, dm := range deviceModels {
		err := storing.DeviceModels.Set(dm)
		if handleError(err, w) {
			return
		}
	}
}
