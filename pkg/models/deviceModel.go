package models

import "strings"

type (
	// DeviceModel is the device capability model decorated with and ID/Name to be used in simulation
	DeviceModel struct {
		ID              string                   `json:"id"`
		Name            string                   `json:"name"`
		CapabilityModel []map[string]interface{} `json:"capabilityModel"`
	}

	// TelemetryType represents a telemetry capability of a device
	TelemetryType struct {
		ID     string
		Name   string
		Schema string
	}

	// PropertyType represents a property capability of a device
	PropertyType struct {
		ID       string
		Name     string
		Schema   string
		Writable bool
	}

	// CommandType represents a command capability of a device
	CommandType struct {
		ID     string
		Name   string
		IsSync bool
	}

	// Component represents a component in the Device Capability Model
	Component struct {
		ComponentID   string
		ComponentType string // FIX THIS LATER!!!
		ComponentName string

		Telemetry  []*TelemetryType
		Properties []*PropertyType
		Commands   []*CommandType
	}

	// DeviceCapabilityModel represents the model of a device
	DeviceCapabilityModel struct {
		Components []*Component
	}
)

// ParseDeviceCapabilityModel parses the DCM from the model store
func (d *DeviceModel) ParseDeviceCapabilityModel() *DeviceCapabilityModel {
	var dcm DeviceCapabilityModel
	var defaultComponent Component
	defaultComponent.ComponentID = "Default"
	defaultComponent.ComponentName = "Default"
	compMap := make(map[string]string)

	for _, component := range d.CapabilityModel {
		var ct Component
		ct.ComponentID = component["@id"].(string)
		ct.ComponentType = component["@type"].(string)
		name, ok := component["displayName"].(map[string]interface{})
		if ok {
			ct.ComponentName = name["en"].(string)
		}

		contents, ok := component["contents"].([]interface{})
		if ok {
			for _, content := range contents {
				var id, typ, name, schema string
				var writable, isSync bool
				for contName, contVal := range content.(map[string]interface{}) {
					if strings.ToLower(contName) == "@type" {
						typ, ok = contVal.(string)
						if !ok {
							types := contVal.([]interface{})
							for _, t := range types {
								tempType := strings.ToLower(t.(string))
								if tempType == "telemetry" {
									typ = tempType
								}
							}
						}
					} else if strings.ToLower(contName) == "@id" {
						id = contVal.(string)
					} else if strings.ToLower(contName) == "name" {
						name = contVal.(string)
					} else if strings.ToLower(contName) == "schema" {
						schema = contVal.(string)
					} else if strings.ToLower(contName) == "writable" {
						writable = contVal.(bool)
					} else if strings.ToLower(contName) == "commandtype" {
						if strings.ToLower(contVal.(string)) == "synchronous" {
							isSync = true
						}
					}
				}
				if strings.ToLower(typ) == "telemetry" {
					ct.Telemetry = append(ct.Telemetry, &TelemetryType{
						ID:     id,
						Name:   name,
						Schema: schema,
					})
				} else if strings.ToLower(typ) == "component" {
					compMap[schema] = name
				} else if strings.ToLower(typ) == "property" {
					ct.Properties = append(ct.Properties, &PropertyType{
						ID:       id,
						Name:     name,
						Schema:   schema,
						Writable: writable,
					})
				} else if strings.ToLower(typ) == "command" {
					ct.Commands = append(ct.Commands, &CommandType{
						ID:     id,
						Name:   name,
						IsSync: isSync,
					})
				}
			}
		}
		dcm.Components = append(dcm.Components, &ct)
	}

	// add default component telemetry
	if len(defaultComponent.Telemetry) > 0 {
		dcm.Components = append(dcm.Components, &defaultComponent)
	}

	// fill in the component names
	for _, comp := range dcm.Components {
		compName, ok := compMap[comp.ComponentID]
		if ok {
			comp.ComponentName = compName
		}
	}
	return &dcm
}
