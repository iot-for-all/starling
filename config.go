package main

import (
	"github.com/iot-for-all/starling/pkg/serving"
	"github.com/iot-for-all/starling/pkg/simulating"
	"github.com/iot-for-all/starling/pkg/storing"
)

type (
	Config struct {
		LogLevel string `yaml:"logLevel" json:"logLevel"` // logging level for the application
		LogsDir  string `yaml:"logsDir" json:"logsDir"`   // directory into which logs are written
	}

	config struct {
		Logger     Config            `yaml:"Logger" json:"Logger"`
		Data       storing.Config    `yaml:"Data" json:"Data"`
		HTTP       serving.Config    `yaml:"HTTP" json:"HTTP"`
		Simulation simulating.Config `yaml:"Simulation" json:"Simulation"`
	}
)

func newConfig() *config {
	return &config{
		Logger: Config{
			LogLevel: "Debug",
			LogsDir:  "./logs",
		},
		Data: storing.Config{
			DataDirectory: "./",
		},
		HTTP: serving.Config{
			AdminPort:   6001,
			MetricsPort: 6002,
		},
		Simulation: simulating.Config{
			ConnectionTimeout:          10000,
			TelemetryTimeout:           10000,
			TwinUpdateTimeout:          10000,
			CommandTimeout:             10000,
			RegistrationAttemptTimeout: 30000,
			MaxConcurrentConnections:   10,
			MaxConcurrentTwinUpdates:   10,
			MaxConcurrentRegistrations: 10,
			MaxConcurrentDeletes:       10,
			MaxRegistrationAttempts:    10,
			EnableTelemetry:            true,
			EnableReportedProps:        false,
			EnableTwinUpdateAcks:       false,
			EnableCommandAcks:          false,
		},
	}
}
