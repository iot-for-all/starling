package serving

import (
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/controlling"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
)

var (
	controller *controlling.Controller
	config     *Config

	//go:embed static
	embeddedFiles embed.FS
)

// StartAdmin starts serving administration API requests.
func StartAdmin(cfg *Config, ctrl *controlling.Controller) {
	config = cfg
	controller = ctrl

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/simulation", listSimulations).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation", upsertSimulation).Methods(http.MethodPut)
	router.HandleFunc("/api/simulation/{id}", getSimulation).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation/{id}", deleteSimulation).Methods(http.MethodDelete)
	router.HandleFunc("/api/simulation/{id}/start", startSimulation).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/stop", stopSimulation).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/provision/{modelId}/{numDevices}", provisionDevices).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/provision/{modelId}/{numDevices}", deleteDevices).Methods(http.MethodDelete)
	router.HandleFunc("/api/simulation/{id}/deviceConfig", listDeviceConfigs).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation/{id}/deviceConfig", upsertDeviceConfig).Methods(http.MethodPut)
	router.HandleFunc("/api/simulation/{id}/deviceConfig/{configId}", getDeviceConfig).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation/{id}/deviceConfig/{configId}", deleteDeviceConfig).Methods(http.MethodDelete)

	router.HandleFunc("/api/target", listTargets).Methods(http.MethodGet)
	router.HandleFunc("/api/target", upsertTarget).Methods(http.MethodPut)
	router.HandleFunc("/api/target/{id}", getTarget).Methods(http.MethodGet)
	router.HandleFunc("/api/target/{id}", deleteTarget).Methods(http.MethodDelete)
	router.HandleFunc("/api/target/{id}/device", listTargetDevices).Methods(http.MethodGet)
	router.HandleFunc("/api/target/{id}/device/{deviceId}", getTargetDevice).Methods(http.MethodGet)
	router.HandleFunc("/api/target/{id}/device", upsertTargetDevice).Methods(http.MethodPut)
	router.HandleFunc("/api/target/{id}/device", deleteAllTargetDevices).Methods(http.MethodDelete)
	router.HandleFunc("/api/target/{id}/device/{deviceId}", deleteTargetDevice).Methods(http.MethodDelete)
	router.HandleFunc("/api/target/{id}/models", getTargetModels).Methods(http.MethodGet)
	router.HandleFunc("/api/target/{id}/models", upsertTargetModels).Methods(http.MethodPut)
	router.HandleFunc("/api/target/{id}/models", deleteTargetModels).Methods(http.MethodDelete)

	router.HandleFunc("/api/model", listDeviceModels).Methods(http.MethodGet)
	router.HandleFunc("/api/model", upsertDeviceModel).Methods(http.MethodPut)
	router.HandleFunc("/api/model/{id}", getDeviceModel).Methods(http.MethodGet)
	router.HandleFunc("/api/model/{id}", deleteDeviceModel).Methods(http.MethodDelete)

	// Serve Starling UX
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(getFileSystem())))

	log.Info().Msgf("serving admin requests at http://localhost:%d/api", cfg.AdminPort)
	log.Info().Msgf("service starling ux at http://localhost:%d", cfg.AdminPort)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", cfg.AdminPort), router)
}

func getFileSystem() http.FileSystem {
	fileSystem, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fileSystem)
}

// handleError log the error and return http error
func handleError(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
