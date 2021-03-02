package main

import (
	"context"
	"fmt"
	"github.com/iot-for-all/starling/pkg/controlling"
	"github.com/iot-for-all/starling/pkg/serving"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
)

func main() {
	// handle process exit gracefully
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer func() {
		// Close the os signal channel to prevent any leak.
		signal.Stop(sig)
	}()

	// load configuration and initialize logger
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to initialize configuration. %w", err))
	}
	initLogger(cfg)

	// initialize database
	err = storing.Open(&cfg.Data)
	if err != nil {
		panic(fmt.Errorf("failed to open the database. %w", err))
	}

	// Initialize the controller.
	controller := controlling.NewController(ctx, &cfg.Simulation)

	// StartSimulation the admin and metrics http endpoints
	go serving.StartAdmin(&cfg.HTTP, controller)
	go serving.StartMetrics(&cfg.HTTP)

	// Wait signal / cancellation
	<-sig

	cancel() // todo: Wait for simulator to completely shut down.
	_ = storing.Close()
}

// loadConfig loads the configuration file
func loadConfig() (*config, error) {
	colorReset := "\033[0m"
	//colorRed := "\033[31m"
	colorGreen := "\033[32m"
	//colorYellow := "\033[33m"
	colorBlue := "\033[34m"
	//colorPurple := "\033[35m"
	//colorCyan := "\033[36m"
	//colorWhite := "\033[37m"
	fmt.Printf(string(colorGreen))
	fmt.Printf(`
   _____ __             ___
  / ___// /_____ ______/ (_)___  ____
  \__ \/ __/ __ \/ ___/ / / __ \/ __ \
 ___/ / /_/ /_/ / /  / / / / / / /_/ /
 ____/\__/\__,_/_/  /_/_/_/ /_/\__, /
                              /____/
`)
	fmt.Printf(string(colorBlue))
	fmt.Printf("     IOT CENTRAL DEVICE SIMULATOR\n")
	fmt.Printf(string(colorReset))

	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	cfgFile := flag.StringP("config", "c", "", "Configuration file to load")
	flag.Parse()

	if *cfgFile != "" {
		viper.SetConfigFile(*cfgFile)
	} else {
		viper.AddConfigPath(home)
		viper.SetConfigName(".starling.yaml")
		viper.SetConfigType("yml")
	}

	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Print(`Add a configuration file ($HOME/.starling.yaml) with the file contents below:

HTTP:
    adminPort: 6001                     # Port number of the administrative service.
    metricsPort: 6002                   # Port number for Prometheus service to scrape.
Simulation:
    connectionTimeout: 30000            # Connection timeout in milli seconds.
    telemetryTimeout: 30000             # Telemetry send timeout in milli seconds.
    twinUpdateTimeout: 30000            # Twin update timeout in milli seconds.
    commandTimeout: 30000               # Command ack timeout in milli seconds.
    registrationAttemptTimeout: 30000   # Device registration timeout in milli seconds.
    maxConcurrentConnections: 100       # Maximum number of concurrent connections to send telemetry per simulation.
    maxConcurrentTwinUpdates: 10        # Maximum number of concurrent twin updates per simulation.
    maxConcurrentRegistrations: 10      # Maximum number of concurrent device registrations (DPS calls).
    maxConcurrentDeletes: 10            # Maximum number of concurrent device deletes.
    maxRegistrationAttempts: 10         # Maximum number of device registration attempts.
    enableTelemetry: true               # Enable device telemetry sends across all simulations.
    enableReportedProps: true           # Enable device reported property sends across all simulations.
    enableTwinUpdateAcks: true          # Enable device twin (desired property) update acknowledgement across all simulations.
    enableCommandAcks: true             # Enable device command (direct method, C2D) acknowledgement across all simulations.
Data:
    dataDirectory: "."                  # Directory used for storing Simulation data.
Logger:
    logLevel: debug				        # Logging legel for the logger. Available logging levels are - panic, fatal, error, warn, info, debug, trace.

\n`)
			return nil, err
		}
	}

	cfg := newConfig()
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if cfg.Data.DataDirectory == "" {
		cfg.Data.DataDirectory = fmt.Sprintf("%s/.starling", home)
	}

	//fmt.Printf("loaded configuration from %s\n", viper.ConfigFileUsed())
	return cfg, nil
}

// initLogger initializes the logger with output format
func initLogger(cfg *config) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	switch strings.ToLower(cfg.Logger.LogLevel) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
