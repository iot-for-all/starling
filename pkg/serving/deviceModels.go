package serving

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"io/ioutil"
	"net/http"
)

// listDeviceModels lists all models.
func listDeviceModels(w http.ResponseWriter, _ *http.Request) {
	items, err := storing.DeviceModels.List()
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(items)
	handleError(err, w)
}

// getDeviceModel gets an existing model by its id.
func getDeviceModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	model, err := storing.DeviceModels.Get(id)
	if handleError(err, w) {
		return
	}

	if model == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(model)
	handleError(err, w)
}

// upsertDeviceModel adds a new or updates an existing device model.
func upsertDeviceModel(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	var model models.DeviceModel
	err = json.Unmarshal(req, &model)
	if handleError(err, w) {
		return
	}

	upsertDeviceModelInternal(w, r, model)
}

// deleteDeviceModel deletes an existing device model.
func deleteDeviceModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := storing.DeviceModels.Delete(id)
	handleError(err, w)
}

func upsertDeviceModelInternal(w http.ResponseWriter, r *http.Request, model models.DeviceModel) {
	err := storing.DeviceModels.Set(&model)
	if handleError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&model)
	handleError(err, w)
}
