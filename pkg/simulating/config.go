package simulating

type Config struct {
	ConnectionTimeout          int  `yaml:"connectionTimeout" json:"connectionTimeout"`
	TelemetryTimeout           int  `yaml:"telemetryTimeout" json:"telemetryTimeout"`
	TwinUpdateTimeout          int  `yaml:"twinUpdateTimeout" json:"twinUpdateTimeout"`
	CommandTimeout             int  `yaml:"commandTimeout" json:"commandTimeout"`
	RegistrationAttemptTimeout int  `yaml:"registrationAttemptTimeout" json:"registrationAttemptTimeout"`
	MaxConcurrentConnections   int  `yaml:"maxConcurrentConnections" json:"maxConcurrentConnections"`
	MaxConcurrentTwinUpdates   int  `yaml:"maxConcurrentTwinUpdates" json:"maxConcurrentTwinUpdates"`
	MaxConcurrentRegistrations int  `yaml:"maxConcurrentRegistrations" json:"maxConcurrentRegistrations"`
	MaxConcurrentDeletes       int  `yaml:"maxConcurrentDeletes" json:"maxConcurrentDeletes"`
	MaxRegistrationAttempts    int  `yaml:"maxRegistrationAttempts" json:"maxRegistrationAttempts"`
	EnableTelemetry            bool `yaml:"enableTelemetry" json:"enableTelemetry"`
	EnableReportedProps        bool `yaml:"enableReportedProps" json:"enableReportedProps"`
	EnableTwinUpdateAcks       bool `yaml:"enableTwinUpdateAcks" json:"enableTwinUpdateAcks"`
	EnableCommandAcks          bool `yaml:"enableCommandAcks" json:"enableCommandAcks"`
}
