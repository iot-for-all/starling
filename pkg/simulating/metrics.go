package simulating

import "github.com/prometheus/client_golang/prometheus"

var (
	simulatedDeviceGauge         *prometheus.GaugeVec
	connectedDeviceGauge         *prometheus.GaugeVec
	connectedDeviceByModelGauge  *prometheus.GaugeVec
	deviceConnectLatency         *prometheus.HistogramVec
	deviceFailoverTotal          *prometheus.CounterVec
	provisionSuccessTotal        *prometheus.CounterVec
	provisionFailuresTotal       *prometheus.CounterVec
	provisionLatency             *prometheus.HistogramVec
	telemetryBatchSuccessTotal   *prometheus.CounterVec
	telemetryBatchSkippedTotal   *prometheus.CounterVec
	telemetryBatchSendLatency    *prometheus.HistogramVec
	telemetryMessageSuccessTotal *prometheus.CounterVec
	telemetryMessageFailureTotal *prometheus.CounterVec
	telemetryMessageSendLatency  *prometheus.HistogramVec
	telemetrySentBytes           *prometheus.CounterVec
	telemetryDataPointsSentTotal *prometheus.CounterVec
	twinUpdateSuccessTotal       *prometheus.CounterVec
	twinUpdateFailureTotal       *prometheus.CounterVec
	twinUpdateSendLatency        *prometheus.HistogramVec
	reportedPropsSkippedTotal    *prometheus.CounterVec
	reportedPropsSuccessTotal    *prometheus.CounterVec
	reportedPropsFailureTotal    *prometheus.CounterVec
	reportedPropsSendLatency     *prometheus.HistogramVec
	commandsSuccessTotal         *prometheus.CounterVec
)

// init initializes the metrics used in simulation
func init() {
	simulatedDeviceGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "device_total",
			Help:      "Total number of devices simulated",
		},
		[]string{"sim", "target", "model"},
	)

	connectedDeviceGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "connect_total",
			Help:      "Total number of devices connected",
		},
		[]string{"sim", "target", "model", "hub"},
	)

	connectedDeviceByModelGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "connect_total_by_model",
			Help:      "Total number of devices connected",
		},
		[]string{"sim", "target", "model"},
	)

	deviceConnectLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "connect_latency_seconds",
			Help:      "Latency of device connecting to IoT Central",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	deviceFailoverTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "failover_total",
			Help:      "Total devices failed over to a new hub",
		},
		[]string{"sim", "target", "model"})

	provisionSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "provisioning",
			Name:      "success_total",
			Help:      "Total devices successfully provisioned",
		},
		[]string{"sim", "target", "model"})

	provisionFailuresTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "provisioning",
			Name:      "failure_total",
			Help:      "Total devices failed provisioning",
		},
		[]string{"sim", "target", "model"},
	)

	provisionLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "provisioning",
			Name:      "latency_seconds",
			Help:      "Latency of provisioning devices",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	telemetryBatchSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_batches_success_total",
			Help:      "Total telemetry batches sent successfully.",
		},
		[]string{"sim", "target", "model"},
	)

	telemetryBatchSkippedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_batches_skipped_total",
			Help:      "Total telemetry batches skipped.",
		},
		[]string{"sim", "target", "model"},
	)

	telemetryBatchSendLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_batch_send_latency_seconds",
			Help:      "Latency of sending telemetry batch from client to IoT Central",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	telemetryMessageSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_messages_success_total",
			Help:      "Total telemetry messages sent successfully.",
		},
		[]string{"sim", "target", "model"},
	)

	telemetryMessageFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_messages_failure_total",
			Help:      "Total telemetry messages send failures.",
		},
		[]string{"sim", "target", "model", "error"},
	)

	telemetryMessageSendLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_message_send_latency_seconds",
			Help:      "Latency of sending telemetry messages from Starling to IoT Central",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	telemetrySentBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_sent_bytes",
			Help:      "Total telemetry data bytes sent.",
		},
		[]string{"sim", "target", "model"},
	)

	telemetryDataPointsSentTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "telemetry_datapoints_sent_total",
			Help:      "Total telemetry data points sent.",
		},
		[]string{"sim", "target", "model"},
	)

	twinUpdateSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "twin_updates_success_total",
			Help:      "Total twin updates sent successfully to IoT Central.",
		},
		[]string{"sim", "target", "model"},
	)

	twinUpdateFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "twin_updates_failure_total",
			Help:      "Total twin update failures.",
		},
		[]string{"sim", "target", "model", "error"},
	)

	twinUpdateSendLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "twin_update_send_latency_seconds",
			Help:      "Latency of sending twin update from client to IoT Central",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	reportedPropsSkippedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "reported_props_skipped_total",
			Help:      "Total reported property updates skipped.",
		},
		[]string{"sim", "target", "model"},
	)

	reportedPropsSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "reported_props_success_total",
			Help:      "Total reported properties sent successfully to IoT Central.",
		},
		[]string{"sim", "target", "model"},
	)

	reportedPropsFailureTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "reported_props_failure_total",
			Help:      "Total reported properties send failures.",
		},
		[]string{"sim", "target", "model", "error"},
	)

	reportedPropsSendLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "reported_props_send_latency_seconds",
			Help:      "Latency of sending reported properties from client to IoT Central",
			Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20, 30, 60, 120, 240, 480, 960},
		},
		[]string{"sim", "target", "model"},
	)

	commandsSuccessTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "starling",
			Subsystem: "simulating",
			Name:      "commands_success_total",
			Help:      "Total successful commands received.",
		},
		[]string{"sim", "target", "model"},
	)

	prometheus.MustRegister(
		simulatedDeviceGauge,
		deviceConnectLatency,
		connectedDeviceGauge,
		connectedDeviceByModelGauge,
		deviceFailoverTotal,
		provisionSuccessTotal,
		provisionFailuresTotal,
		provisionLatency,
		telemetryBatchSuccessTotal,
		telemetryBatchSkippedTotal,
		telemetryBatchSendLatency,
		telemetryMessageSuccessTotal,
		telemetryMessageFailureTotal,
		telemetryMessageSendLatency,
		telemetrySentBytes,
		telemetryDataPointsSentTotal,
		twinUpdateSuccessTotal,
		twinUpdateFailureTotal,
		twinUpdateSendLatency,
		reportedPropsSkippedTotal,
		reportedPropsSuccessTotal,
		reportedPropsFailureTotal,
		reportedPropsSendLatency,
		commandsSuccessTotal,
	)
}
