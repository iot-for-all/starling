package serving

// Config containing http configuration
type Config struct {
	AdminPort   int `yaml:"adminPort" json:"adminPort"`     // port number of the administration API server
	MetricsPort int `yaml:"metricsPort" json:"metricsPort"` // port number of prometheus metrics server
}
