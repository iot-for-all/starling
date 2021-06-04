package serving

import (
	"embed"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iot-for-all/starling/pkg/config"
	"github.com/iot-for-all/starling/pkg/controlling"
	"github.com/rs/zerolog/log"
	"io/fs"
	"net/http"
	"os"
	"path"
)

var (
	controller   *controlling.Controller
	globalConfig *config.GlobalConfig

	//go:embed static
	embeddedFiles embed.FS
)

// StartAdmin starts serving administration API requests.
func StartAdmin(globalCfg *config.GlobalConfig, ctrl *controlling.Controller) {
	globalConfig = globalCfg
	controller = ctrl

	router := mux.NewRouter().StrictSlash(true)

	// Service API
	router.HandleFunc("/api/simulation", listSimulations).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation", upsertSimulation).Methods(http.MethodPut)
	router.HandleFunc("/api/simulation/{id}", getSimulation).Methods(http.MethodGet)
	router.HandleFunc("/api/simulation/{id}", deleteSimulation).Methods(http.MethodDelete)
	router.HandleFunc("/api/simulation/{id}/start", startSimulation).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/stop", stopSimulation).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/provision/{modelId}/{numDevices}", provisionDevices).Methods(http.MethodPost)
	router.HandleFunc("/api/simulation/{id}/provision", deleteAllDevices).Methods(http.MethodDelete)
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

	// WEB API
	router.HandleFunc("/webapi/model", webAPIListDeviceModels).Methods(http.MethodGet)
	router.HandleFunc("/webapi/model/{id}", webAPIGetDeviceModel).Methods(http.MethodGet)
	router.HandleFunc("/webapi/model", webAPIAddDeviceModel).Methods(http.MethodPost)
	router.HandleFunc("/webapi/model", webAPIUpdateDeviceModel).Methods(http.MethodPut)
	router.HandleFunc("/webapi/model/{id}", webAPIDeleteDeviceModel).Methods(http.MethodDelete)

	router.HandleFunc("/webapi/target", webAPIListTargets).Methods(http.MethodGet)
	router.HandleFunc("/webapi/target/{id}", webAPIGetTarget).Methods(http.MethodGet)
	router.HandleFunc("/webapi/target", webAPIAddTarget).Methods(http.MethodPost)
	router.HandleFunc("/webapi/target", webAPIUpdateTarget).Methods(http.MethodPut)
	router.HandleFunc("/webapi/target/{id}", webAPIDeleteTarget).Methods(http.MethodDelete)
	router.HandleFunc("/webapi/target/{id}/import", webAPIImportModelsFromTarget).Methods(http.MethodPost)

	router.HandleFunc("/webapi/simulation", webAPIListSimulations).Methods(http.MethodGet)
	router.HandleFunc("/webapi/simulation/{id}", webAPIGetSimulation).Methods(http.MethodGet)
	router.HandleFunc("/webapi/simulation", webAPIAddSimulation).Methods(http.MethodPost)
	router.HandleFunc("/webapi/simulation", webAPIUpdateSimulation).Methods(http.MethodPut)
	router.HandleFunc("/webapi/simulation/{id}", webAPIDeleteSimulation).Methods(http.MethodDelete)
	router.HandleFunc("/webapi/simulation/{id}/provision", webAPIProvisionDevices).Methods(http.MethodPost)
	router.HandleFunc("/webapi/simulation/{id}/start", webAPIStartSimulation).Methods(http.MethodPost)
	router.HandleFunc("/webapi/simulation/{id}/stop", webAPIStopSimulation).Methods(http.MethodPost)
	router.HandleFunc("/webapi/simulation/{id}/export", webAPIExportSimulation).Methods(http.MethodGet)

	router.HandleFunc("/webapi/config", webAPIGetConfig).Methods(http.MethodGet)
	router.HandleFunc("/webapi/config", webAPIUpdateConfig).Methods(http.MethodPut)
	router.HandleFunc("/webapi/config/metricsStatus", webAPIMetricsStatus).Methods(http.MethodGet)

	// Serve Starling UX
	handler := http.FileServer(getFileSystem())
	router.PathPrefix("/").Handler(http.StripPrefix("/", handler))
	//router.PathPrefix("/sim").Handler(http.StripPrefix("/sim", handler))
	//handler := AssetHandler("/", "static")
	//router.HandleFunc("/", handler.ServeHTTP).Methods(http.MethodGet)
	//router.HandleFunc("/*", handler.ServeHTTP).Methods(http.MethodGet)

	log.Info().Msgf("serving admin requests at http://localhost:%d/api", globalConfig.HTTP.AdminPort)
	log.Info().Msgf("serving UX at http://localhost:%d", globalConfig.HTTP.AdminPort)

	// handle CORS
	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
	originsOK := handlers.AllowedOrigins([]string{"*"})

	_ = http.ListenAndServe(fmt.Sprintf(":%d", globalConfig.HTTP.AdminPort), handlers.CORS(headersOK, methodsOK, originsOK)(router))
}

func getFileSystem() http.FileSystem {
	fileSystem, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fileSystem)
}

// fsFunc is short-hand for constructing a http.FileSystem
// implementation
type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

// AssetHandler returns an http.Handler that will serve files from
// the Assets embed.FS.  When locating a file, it will strip the given
// prefix from the request and prepend the root to the filesystem
// lookup: typical prefix might be /web/, and root would be build.
func AssetHandler(prefix, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		log.Debug().Str("assetPath", assetPath).Msg("got request")
		// If we can't find the asset, return the default index.html
		// content
		f, err := embeddedFiles.Open(assetPath)
		if os.IsNotExist(err) {
			log.Debug().Str("assetPath", assetPath).Msg("got request redirected to index.html")
			return embeddedFiles.Open("index.html")
		}

		// Otherwise assume this is a legitimate request routed
		// correctly
		return f, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}

// handleError log the error and return http error
func handleError(err error, w http.ResponseWriter) bool {
	if err != nil {
		log.Error().Err(err).Msg("error encountered while processing request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
