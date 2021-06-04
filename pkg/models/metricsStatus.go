package models

type MetricsStatus struct {
	GrafanaServer    bool `json:"grafanaServer"`
	PrometheusServer bool `json:"prometheusServer"`
}
