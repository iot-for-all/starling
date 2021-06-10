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
	"strings"
)

// webAPIListDeviceModels lists all models.
func webAPIListDeviceModels(w http.ResponseWriter, r *http.Request) {
	listDeviceModels(w, r)
}

// webAPIGetDeviceModel gets an existing model by its id.
func webAPIGetDeviceModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if len(id) == 0 {
		msg := fmt.Sprintf("model id is required for GET device model call.")
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	getDeviceModel(w, r)
}

// webAPIUpdateDeviceModel adds a new or updates an existing device model.
func webAPIAddDeviceModel(w http.ResponseWriter, r *http.Request) {
	// check if an existing device model exists with the same ID
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var model models.DeviceModel
	err = json.Unmarshal(req, &model)
	if handleError(err, w) {
		return
	}

	// generate ID if needed
	if len(model.ID) == 0 {
		model.ID, err = generateModelID(model.Name)
		if handleError(err, w) {
			return
		}
	}

	dm, err := storing.DeviceModels.Get(model.ID)
	if handleError(err, w) {
		return
	}

	if dm != nil {
		msg := fmt.Sprintf("Another device model exists with Model ID '%s'. Try again with another ID.", model.ID)
		log.Error().Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	upsertDeviceModelInternal(w, r, model)

	// add this device model to all targets
	targets, err := storing.Targets.List()
	if handleError(err, w) {
		return
	}

	for _, target := range targets {
		m, err := storing.TargetModels.Get(target.ID)
		if handleError(err, w) {
			return
		}

		if m != nil {
			m.Models = append(m.Models, model.ID)

			err = storing.TargetModels.Set(m)
			if handleError(err, w) {
				return
			}
		}
	}
}

// webAPIUpdateDeviceModel updates an existing device model.
func webAPIUpdateDeviceModel(w http.ResponseWriter, r *http.Request) {
	upsertDeviceModel(w, r)
}

// webAPIDeleteDeviceModel deletes an existing device model.
func webAPIDeleteDeviceModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// check if there are any simulations using this model
	sims, err := storing.Simulations.List()
	if handleError(err, w) {
		return
	}

	for _, sim := range sims {
		deviceConfigs, err := storing.DeviceConfigs.List(sim.ID)
		if handleError(err, w) {
			return
		}

		for _, dc := range deviceConfigs {
			// TODO: Check for provisioned and sim device counts
			if dc.ModelID == id {
				msg := fmt.Sprintf("device model '%s' cannot be deleted as there are simulations using this model. Delete the simulation '%s' and try again.", id, sim.Name)
				log.Error().Msg(msg)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
		}
	}

	deleteDeviceModel(w, r)

	// delete this device model from all targets
	targets, err := storing.Targets.List()
	if handleError(err, w) {
		return
	}

	for _, target := range targets {
		err = storing.TargetModels.Delete(target.ID)
		if handleError(err, w) {
			return
		}
	}
}

func generateModelID(name string) (string, error) {
	name = strings.ToLower(name)
	allowedChars := "abcdefghijklmnopqrstuvwxyz0123456789"
	id := ""
	for _, ch := range name {
		char := string(ch)
		if strings.Index(allowedChars, char) >= 0 {
			id += char
		}
	}
	if len(id) == 0 {
		id = "model"
	}

	count := 0
	newId := id
	for {
		// check to see if there are existing models with same id
		model, err := storing.DeviceModels.Get(newId)
		if err != nil {
			return id, err
		}

		if model == nil {
			// no model exists with that id, it is good to use
			return newId, nil
		}

		count++
		newId = fmt.Sprintf("%s%d", id, count)
	}
}
