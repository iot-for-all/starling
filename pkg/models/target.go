package models

type SimulationTargetType string

type (
	// SimulationTarget specifies the target of a simulation
	SimulationTarget struct {
		ID              string `json:"id"`              // user supplied identifier of a target.
		Name            string `json:"name"`            // display name of the target.
		ProvisioningURL string `json:"provisioningUrl"` // DPS provisioning URL.
		IDScope         string `json:"idScope"`         // the id scope of the provisioning endpoint.
		MasterKey       string `json:"masterKey"`       // the master SAS key of the provisioning endpoint.
		AppUrl          string `json:"appUrl"`          // Central app URL
		AppToken        string `json:"appToken"`        // Central app token for API access
	}

	// SimulationTargetModels specifies the models configured for a simulation target.
	SimulationTargetModels struct {
		TargetID string   `json:"targetId"` // identifier of a target.
		Models   []string `json:"models"`   // list of model IDs available for this target.
	}

	// SimulationTargetDevice cached copy of a device connection string  in a target.
	SimulationTargetDevice struct {
		TargetID         string `json:"targetId"`         // identifier of a target.
		DeviceID         string `json:"deviceId"`         // device identifier in the target.
		ConnectionString string `json:"connectionString"` // IoT Hub connection string for the device.
	}
)
