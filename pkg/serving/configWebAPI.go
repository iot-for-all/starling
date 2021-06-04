package serving

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iot-for-all/starling/pkg/config"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// webAPIGetConfig get the current configuration.
func webAPIGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(globalConfig)
	handleError(err, w)
}

// webAPIUpdateConfig update current configuration.
func webAPIUpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := ioutil.ReadAll(r.Body)
	if handleError(err, w) {
		return
	}

	// read updated config from request
	var cfg config.GlobalConfig
	err = json.Unmarshal(req, &cfg)
	if handleError(err, w) {
		return
	}

	// update config
	globalConfig.Data = cfg.Data
	globalConfig.HTTP = cfg.HTTP
	globalConfig.Logger = cfg.Logger
	globalConfig.Simulation = cfg.Simulation

	// generate YAML content and write it to the config file
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("error reading current directory: %s\n", err)
	}
	configFileName := path.Join(exeDir, "starling.json")
	content, err := json.MarshalIndent(cfg, "", "  ")
	if handleError(err, w) {
		return
	}
	err = os.WriteFile(configFileName, content, os.ModePerm)
	if handleError(err, w) {
		return
	}

	// write the config back to request
	err = json.NewEncoder(w).Encode(globalConfig)
	handleError(err, w)
}

// webAPIMetricsStatus gets the status of grafana and prometheus servers
func webAPIMetricsStatus(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	metricsStatus := controller.GetMetricsStatus(ctx)

	// write the config back to request
	err := json.NewEncoder(w).Encode(metricsStatus)
	handleError(err, w)
}
