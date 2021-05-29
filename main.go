package main

import (
	"context"
	"fmt"
	"github.com/iot-for-all/starling/pkg/config"
	"github.com/iot-for-all/starling/pkg/controlling"
	"github.com/iot-for-all/starling/pkg/serving"
	"github.com/iot-for-all/starling/pkg/storing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
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
	controller.ResetSimulationStatus()

	// Start the admin and metrics http endpoints
	go serving.StartAdmin(cfg, controller)
	go serving.StartMetrics(&cfg.HTTP)

	// open web browser serving the Starling website
	url := fmt.Sprintf("http://localhost:%d", cfg.HTTP.AdminPort)
	err = openWebBrowser(url)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("failed to open web browser with url %s", url))
	}

	// Wait signal / cancellation
	<-sig

	cancel() // todo: Wait for simulator to completely shut down.
	_ = storing.Close()
}

// loadConfig loads the configuration file
func loadConfig() (*config.GlobalConfig, error) {
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

	cfg := config.NewConfig()

	// if the config file does not exist, write a default config file
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("error reading current directory: %s\n", err)
	}
	configFileName := path.Join(exeDir, "starling.yaml")
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		content, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Printf("error reading current directory: %s\n", err)
		} else {
			err = os.WriteFile(configFileName, content, os.ModePerm)
			if err != nil {
				fmt.Printf("error writing to default config file %s: %s\n", configFileName, err)
			}
		}
	}

	// read the config file
	content, err := os.ReadFile(configFileName)
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		fmt.Printf("error reading config file %s: %s", configFileName, err)
	}

	if cfg.Data.DataDirectory == "" {
		cfg.Data.DataDirectory = fmt.Sprintf("%s/data", exeDir)
	}

	return cfg, nil
}

// initLogger initializes the logger with output format
func initLogger(cfg *config.GlobalConfig) {
	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

	fileLoggingEnabled := false
	if len(cfg.Logger.LogsDir) > 0 {
		fileLoggingEnabled = true
	}
	if fileLoggingEnabled {
		logsDir := cfg.Logger.LogsDir
		if err := os.MkdirAll(logsDir, 0744); err != nil {
			fmt.Printf("can't create log directory, so file logging is disabled, error: %s", err.Error())
		} else {
			fileWriter := &lumberjack.Logger{
				Filename:   path.Join(logsDir, "starling.log"),
				MaxBackups: 3,  // files
				MaxSize:    10, // megabytes
				MaxAge:     30, // days
			}

			writers = append(writers, fileWriter)
			//fmt.Printf("file logging is enabled, logsDir: %s\n", logsDir)
		}
	}
	mw := io.MultiWriter(writers...)

	log.Logger = zerolog.New(mw).With().Timestamp().Logger()
	//log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

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

// openWebBrowser opens the specified URL in the default browser of the user.
func openWebBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
