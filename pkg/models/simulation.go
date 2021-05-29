package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type (
	// DeviceDisconnectBehavior defines the disconnection behavior of the simulated device.
	DeviceDisconnectBehavior string

	// TelemetryFormat defines the format of the telemetry messages sent from the simulated device.
	TelemetryFormat string

	// SimulationStatus specifies the current status of the simulation.
	SimulationStatus string

	// SimulationDeviceConfig defines the device configuration for a simulation.
	SimulationDeviceConfig struct {
		ID          string `json:"id"`          // the id of the configuration
		ModelID     string `json:"modelId"`     // the model to simulate.
		DeviceCount int    `json:"deviceCount"` // the total no. of devices to simulate.
	}

	// Simulation definition.
	Simulation struct {
		ID                    string                   `json:"id"`                       // the id of the simulation.
		Name                  string                   `json:"name"`                     // the name of the simulation.
		TargetID              string                   `json:"targetId"`                 // the id of the target application against which the simulation is running.
		Status                SimulationStatus         `json:"status"`                   // current status of the simulation.
		WaveGroupCount        int                      `json:"waveGroupCount"`           // no. of device groups in a wave in which the total devices are distributed.
		WaveGroupInterval     int                      `json:"waveGroupInterval"`        // interval between wave groups when simulating.
		TelemetryBatchSize    int                      `json:"telemetryBatchSize"`       // batch of telemetry messages that each device will send.
		TelemetryInterval     int                      `json:"telemetryInterval"`        // interval to wait between sending telemetry messages.
		ReportedPropsInterval int                      `json:"reportedPropertyInterval"` // interval to wait between sending reported properties.
		DisconnectBehavior    DeviceDisconnectBehavior `json:"disconnectBehavior"`       // device connection behavior.
		TelemetryFormat       TelemetryFormat          `json:"telemetryFormat"`          // format of telemetry messages.
		LastUpdatedTime       time.Time                `json:"lastUpdatedTime"`          // when the status was last updated
	}

	// SimulationViewDeviceConfig defines the device configuration for a simulation view.
	SimulationViewDeviceConfig struct {
		ID               string `json:"id"`               // the id of the configuration
		ModelID          string `json:"modelId"`          // the model to simulate.
		ProvisionedCount int    `json:"provisionedCount"` // number of provisioned devices.
		SimulatedCount   int    `json:"simulatedCount"`   // number of devices to simulate.
		ConnectedCount   int    `json:"connectedCount"`   // number of devices currently connected.
	}

	SimulationView struct {
		Simulation                              // simulation configuration
		Devices    []SimulationViewDeviceConfig `json:"devices"` // devices configurations
	}
)

const (
	// SimulationStatusReady specifies that the simulation is ready to run or provision devices.
	SimulationStatusReady SimulationStatus = "ready"
	// SimulationStatusRunning specifies that the simulation is running.
	SimulationStatusRunning SimulationStatus = "running"
	// SimulationStatusProvisioning specifies that the simulation is provisioning devices.
	SimulationStatusProvisioning SimulationStatus = "provisioning"
	// SimulationStatusDeleting specifies that the devices in the simulation are getting deleted.
	SimulationStatusDeleting SimulationStatus = "deleting"

	// DeviceDisconnectNever specifies that the device should never disconnect.
	DeviceDisconnectNever DeviceDisconnectBehavior = "never"
	// DeviceDisconnectAfterTelemetrySend specifies that the device should disconnect after sending telemetry.
	DeviceDisconnectAfterTelemetrySend DeviceDisconnectBehavior = "telemetry"

	// TelemetryFormatDefault specifies that the device sends telemetry in default JSON format.
	TelemetryFormatDefault TelemetryFormat = "default"
	// TelemetryFormatOpcua specifies that the device sends telemetry in opcua JSON format.
	TelemetryFormatOpcua TelemetryFormat = "opcua"
)

// UnmarshalJSON handles the un-marshalling of simulation status
func (status *SimulationStatus) UnmarshalJSON(b []byte) error {
	var p string
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p == "" {
		return nil
	}

	// backward compat
	if p == "created" || p == "stopped" {
		p = "ready"
	}

	s := SimulationStatus(p)
	switch s {
	case SimulationStatusReady,
		SimulationStatusRunning,
		SimulationStatusProvisioning,
		SimulationStatusDeleting:
		*status = s
		return nil
	default:
		return fmt.Errorf("invalid simulation status type %s", p)
	}
}

// UnmarshalJSON handles the un-marshalling of device disconnect behavior.
func (d *DeviceDisconnectBehavior) UnmarshalJSON(b []byte) error {
	var p string
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p == "" {
		return nil
	}

	s := DeviceDisconnectBehavior(p)
	switch s {
	case DeviceDisconnectNever,
		DeviceDisconnectAfterTelemetrySend:
		*d = s
		return nil
	default:
		return fmt.Errorf("invalid device disconnect type %s", p)
	}
}

// UnmarshalJSON handles the un-marshalling of telemetry format
func (tf *TelemetryFormat) UnmarshalJSON(b []byte) error {
	var p string
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	if p == "" {
		return nil
	}

	s := TelemetryFormat(p)
	switch s {
	case TelemetryFormatDefault,
		TelemetryFormatOpcua:
		*tf = s
		return nil
	default:
		return fmt.Errorf("invalid telemetry format type %s", p)
	}
}
