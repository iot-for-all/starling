package serving

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/iot-for-all/starling/pkg/storing"
	"strings"
	"time"
)

type exporter struct {
	sim *models.Simulation
}

func (e *exporter) exportSimulation() ([]byte, error) {
	builder := strings.Builder{}

	e.exportFileHeader(&builder)

	target, err := storing.Targets.Get(e.sim.TargetID)
	if err != nil {
		return nil, err
	}

	err = e.exportTarget(&builder, target)
	if err != nil {
		return nil, err
	}

	deviceConfigs, err := storing.DeviceConfigs.List(e.sim.ID)
	if err != nil {
		return nil, err
	}

	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Add device models.\n")
	builder.WriteString("## These device models can be used across all simulations.\n")
	modelNames := make([]string, len(deviceConfigs))
	for i, dc := range deviceConfigs {
		model, err := storing.DeviceModels.Get(dc.ModelID)
		if err != nil {
			return nil, err
		}

		modelNames[i] = dc.ModelID
		err = e.exportModel(&builder, model)
		if err != nil {
			return nil, err
		}
	}

	err = e.exportModelBinding(&builder, target.ID, modelNames)
	if err != nil {
		return nil, err
	}

	err = e.exportSimulationDetails(&builder, e.sim)
	if err != nil {
		return nil, err
	}

	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Configure number of devices for the simulation.\n")
	builder.WriteString("## You can have multiple device types per simulation by adding a several\n")
	builder.WriteString("## device configs just like the one below.\n")
	for _, dc := range deviceConfigs {
		err = e.exportSimulationDeviceConfig(&builder, e.sim, dc)
		if err != nil {
			return nil, err
		}
	}
	e.exportStartStopSim(&builder, e.sim)

	return []byte(builder.String()), nil
}

func (e *exporter) exportFileHeader(builder *strings.Builder) {
	builder.WriteString("#!/bin/sh\n")
	builder.WriteString("#####################################################################################\n")
	builder.WriteString(fmt.Sprintf("## Generated from Startling at: %s\n", time.Now().Format(time.RFC3339)))
	builder.WriteString("##\n")
	builder.WriteString("## This shell script can be used to seed all the data for the simulation.\n")
	builder.WriteString("## Start starling server and run this script to seed data.\n")
	builder.WriteString("##\n")
	builder.WriteString("## Start/Stop simulation commands are commented out at the bottom opf the file.\n")
	builder.WriteString("#####################################################################################\n\n")
	builder.WriteString("## Change this parameter based on your setup\n")
	builder.WriteString("BASE_URL=\"http://localhost:6001/api\"       # Starling API endpoint\n\n")
}

func (e *exporter) exportTarget(builder *strings.Builder, target *models.SimulationTarget) error {
	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Add an IoT Central application in which the simulated devices will be created.\n")
	builder.WriteString("## You can add multiple applications similar to the one below.\n")
	builder.WriteString(fmt.Sprintf("## Target Application: '%s'\n", target.Name))
	builder.WriteString("curl --location --request PUT \"$BASE_URL/target\" \\\n")
	builder.WriteString("--header 'Content-Type: application/json' \\\n")
	builder.WriteString("--data-raw '\n")
	content, err := e.marshallToJson(target)

	if err != nil {
		return err
	}
	builder.Write(content)
	builder.WriteString("'\n\n")
	return nil
}

func (e *exporter) exportModel(builder *strings.Builder, model *models.DeviceModel) error {
	builder.WriteString(fmt.Sprintf("## Device model: %s\n", model.Name))
	builder.WriteString("curl --location --request PUT \"$BASE_URL/model\" \\\n")
	builder.WriteString("--header 'Content-Type: application/json' \\\n")
	builder.WriteString("--data-raw '\n")
	content, err := e.marshallToJson(model)
	if err != nil {
		return err
	}
	builder.Write(content)
	builder.WriteString("'\n\n")
	return nil
}

func (e *exporter) exportModelBinding(builder *strings.Builder, targetId string, modelNames []string) error {
	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Add target application to model bindings so that devices of this type\n")
	builder.WriteString("## can be created in simulations.\n")
	builder.WriteString("## You can add multiple device models in the same models array.\n")
	builder.WriteString("##\n")
	builder.WriteString("## NOTE: You need to add this model to IoT Central application yourself\n")
	builder.WriteString("##       before the simulation is started.\n")
	builder.WriteString(fmt.Sprintf("curl --location --request PUT \"$BASE_URL/target/%s/models\" \\\n", targetId))
	builder.WriteString("--header 'Content-Type: application/json' \\\n")
	builder.WriteString("--data-raw '\n")
	content, err := e.marshallToJson(modelNames)
	if err != nil {
		return err
	}
	builder.Write(content)
	builder.WriteString("'\n\n")
	return nil
}

func (e *exporter) exportSimulationDetails(builder *strings.Builder, simulation *models.Simulation) error {
	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Add a Simulation.\n")
	builder.WriteString(fmt.Sprintf("## This simulation is configured to distribute the devices into %d wave group(s).\n", simulation.WaveGroupCount))
	builder.WriteString(fmt.Sprintf("## Each wave group will send data %d second(s) apart.\n", simulation.WaveGroupInterval))
	builder.WriteString(fmt.Sprintf("## Telemetry is sent every %d second(s) and each time it sends a batch of %d message(s).\n", simulation.TelemetryInterval, simulation.TelemetryBatchSize))
	builder.WriteString(fmt.Sprintf("## Reported properties are sent every %d second(s).\n", simulation.ReportedPropsInterval))
	builder.WriteString("## Device disconnect behavior 'never' means that the devices are never disconnected.\n")
	builder.WriteString("## To simulate an occasionally connected device, you can change the disconnectBehavior\n")
	builder.WriteString("## to 'telemetry' to disconnect the device after sending telemetry.\n")
	builder.WriteString(fmt.Sprintf("## Simulation: %s\n", simulation.Name))
	builder.WriteString("curl --location --request PUT \"$BASE_URL/simulation\" \\\n")
	builder.WriteString("--header 'Content-Type: application/json' \\\n")
	builder.WriteString("--data-raw '\n")
	content, err := e.marshallToJson(simulation)
	if err != nil {
		return err
	}
	builder.Write(content)
	builder.WriteString("'\n\n")
	return nil
}

func (e *exporter) exportSimulationDeviceConfig(builder *strings.Builder, simulation *models.Simulation, deviceConfig *models.SimulationDeviceConfig) error {
	builder.WriteString(fmt.Sprintf("## Setup %d %s devices in simulation %s\n", deviceConfig.DeviceCount, deviceConfig.ModelID, simulation.Name))
	builder.WriteString(fmt.Sprintf("curl --location --request PUT \"$BASE_URL/simulation/%s/deviceConfig\" \\\n", simulation.ID))
	builder.WriteString("--header 'Content-Type: application/json' \\\n")
	builder.WriteString("--data-raw '\n")
	content, err := e.marshallToJson(deviceConfig)
	if err != nil {
		return err
	}
	builder.Write(content)
	builder.WriteString("'\n\n")
	return nil
}

func (e *exporter) exportStartStopSim(builder *strings.Builder, simulation *models.Simulation) {
	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Start this simulation.\n")
	builder.WriteString(fmt.Sprintf("## Use this command to start the simulation %s\n", simulation.Name))
	builder.WriteString(fmt.Sprintf("# curl --location --request POST \"$BASE_URL/simulation/%s/start\"\n", simulation.ID))
	builder.WriteString("\n\n")

	builder.WriteString("#####################################################################################\n")
	builder.WriteString("## Stop this simulation.\n")
	builder.WriteString(fmt.Sprintf("## Use this command to stop the simulation %s\n", simulation.Name))
	builder.WriteString(fmt.Sprintf("# curl --location --request POST \"$BASE_URL/simulation/%s/stop\"\n", simulation.ID))
}

// marshallToJson converts the given structure to JSON without HTML escaping strings
func (e *exporter) marshallToJson(obj interface{}) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buff)
	jsonEncoder.SetIndent("", "  ")
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(obj)
	return buff.Bytes(), err
}
