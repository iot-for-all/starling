package config

type (
	LoggerConfig struct {
		LogLevel string `yaml:"logLevel" json:"logLevel"` // logging level for the application
		LogsDir  string `yaml:"logsDir" json:"logsDir"`   // directory into which logs are written
	}

	StoreConfig struct {
		DataDirectory string `yaml:"path" json:"path"`
	}

	HTTPConfig struct {
		AdminPort   int `yaml:"adminPort" json:"adminPort"`     // port number of the administration API server
		MetricsPort int `yaml:"metricsPort" json:"metricsPort"` // port number of prometheus metrics server
	}

	SimulationConfig struct {
		ConnectionTimeout          int          `yaml:"connectionTimeout" json:"connectionTimeout"`
		TelemetryTimeout           int          `yaml:"telemetryTimeout" json:"telemetryTimeout"`
		TwinUpdateTimeout          int          `yaml:"twinUpdateTimeout" json:"twinUpdateTimeout"`
		CommandTimeout             int          `yaml:"commandTimeout" json:"commandTimeout"`
		RegistrationAttemptTimeout int          `yaml:"registrationAttemptTimeout" json:"registrationAttemptTimeout"`
		MaxConcurrentConnections   int          `yaml:"maxConcurrentConnections" json:"maxConcurrentConnections"`
		MaxConcurrentTwinUpdates   int          `yaml:"maxConcurrentTwinUpdates" json:"maxConcurrentTwinUpdates"`
		MaxConcurrentRegistrations int          `yaml:"maxConcurrentRegistrations" json:"maxConcurrentRegistrations"`
		MaxConcurrentDeletes       int          `yaml:"maxConcurrentDeletes" json:"maxConcurrentDeletes"`
		MaxRegistrationAttempts    int          `yaml:"maxRegistrationAttempts" json:"maxRegistrationAttempts"`
		EnableTelemetry            bool         `yaml:"enableTelemetry" json:"enableTelemetry"`
		EnableReportedProps        bool         `yaml:"enableReportedProps" json:"enableReportedProps"`
		EnableTwinUpdateAcks       bool         `yaml:"enableTwinUpdateAcks" json:"enableTwinUpdateAcks"`
		EnableCommandAcks          bool         `yaml:"enableCommandAcks" json:"enableCommandAcks"`
		GeopointData               [][3]float64 `yaml:"geopointData" json:"geopointData"`
	}

	GlobalConfig struct {
		Logger     LoggerConfig     `yaml:"Logger" json:"Logger"`
		Data       StoreConfig      `yaml:"Data" json:"Data"`
		HTTP       HTTPConfig       `yaml:"HTTP" json:"HTTP"`
		Simulation SimulationConfig `yaml:"Simulation" json:"Simulation"`
	}
)

func NewConfig() *GlobalConfig {
	return &GlobalConfig{
		Logger: LoggerConfig{
			LogLevel: "debug",
			LogsDir:  "./logs",
		},
		Data: StoreConfig{
			DataDirectory: "./",
		},
		HTTP: HTTPConfig{
			AdminPort:   6001,
			MetricsPort: 6002,
		},
		Simulation: SimulationConfig{
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
			EnableReportedProps:        true,
			EnableTwinUpdateAcks:       true,
			EnableCommandAcks:          true,
			GeopointData: [][3]float64{
				{47.645804, -122.132337, 0.0},
				{47.644799, -122.132291, 0.0},
				{47.643975, -122.132302, 0.0},
				{47.642746, -122.132366, 0.0},
				{47.641264, -122.132409, 0.0},
				{47.639768, -122.132430, 0.0},
				{47.637844, -122.132393, 0.0},
				{47.635111, -122.132479, 0.0},
				{47.633202, -122.132382, 0.0},
				{47.633354, -122.131191, 0.0},
				{47.634540, -122.129163, 0.0},
				{47.636325, -122.126081, 0.0},
				{47.638046, -122.123120, 0.0},
				{47.641111, -122.119204, 0.0},
				{47.644017, -122.115642, 0.0},
				{47.645990, -122.114258, 0.0},
				{47.646069, -122.117938, 0.0},
				{47.646069, -122.120921, 0.0},
				{47.646105, -122.125888, 0.0},
				{47.646047, -122.129568, 0.0},
				{47.646069, -122.132164, 0.0},
			},
		},
	}
}
