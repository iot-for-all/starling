package simulating

import (
	"encoding/json"
	"fmt"
	"github.com/amenzhinsky/iothub/common"
	"github.com/amenzhinsky/iothub/iotdevice"
	"github.com/hashicorp/go-uuid"
	"github.com/iot-for-all/starling/pkg/models"
	"math/rand"
	"strings"
	"time"
)

type (
	// DataGenerator generates telemetry messages and reported property updates based on the device capability model.
	DataGenerator struct {
		CapabilityModel *models.DeviceCapabilityModel // the capability model of the device.
		nextGeoPoint    int                           // geo point to be used next from the geopointRoute
	}
)

var (
	geopointRoute = [][2]float64{
		{47.645804, -122.132337},
		{47.644799, -122.132291},
		{47.643975, -122.132302},
		{47.642746, -122.132366},
		{47.641264, -122.132409},
		{47.639768, -122.132430},
		{47.637844, -122.132393},
		{47.635111, -122.132479},
		{47.633202, -122.132382},
		{47.633354, -122.131191},
		{47.634540, -122.129163},
		{47.636325, -122.126081},
		{47.638046, -122.123120},
		{47.641111, -122.119204},
		{47.644017, -122.115642},
		{47.645990, -122.114258},
		{47.646069, -122.117938},
		{47.646069, -122.120921},
		{47.646105, -122.125888},
		{47.646047, -122.129568},
		{47.646069, -122.132164},
	}
)

// GenerateTelemetryMessage generate a telemetry messages based on the device capability model.
func (d *DataGenerator) GenerateTelemetryMessage(device *device, creationTime time.Time) ([]*telemetryMessage, error) {
	// TODO: Handle components

	telemetryMessages := make([]*telemetryMessage, 1)

	dataPointCount := 0
	var body []byte
	var err error
	if device.simulation.TelemetryFormat == models.TelemetryFormatOpcua {
		// OPCUA device sending JSON payload
		msgGuid, _ := uuid.GenerateUUID()
		payload := make(map[string]interface{})
		msgList := make([]map[string]interface{}, 1)
		device.telemetrySequenceNumber++
		msgList[0] = map[string]interface{}{
			"DataSetWriterId": fmt.Sprintf("%s-%s", device.deviceID, msgGuid),
			"MetaDataVersion": map[string]interface{}{
				"MajorVersion": 1,
				"MinorVersion": 0,
			},
			"SequenceNumber": device.telemetrySequenceNumber, //  rand.Intn(100000),
			"Status":         nil,
			"Timestamp":      d.getDateTime(),
			"Payload":        payload,
		}

		eventId, _ := uuid.GenerateUUID()
		telemetryValues := map[string]interface{}{
			"DataSetClassId":     nil,
			"DataSetWriterGroup": device.deviceID,
			"EventId":            eventId,
			"MessageId":          d.getString(5),
			"MessageType":        "ua-data",
			"PublisherId":        "Standalone_IIOTEdgeServer_opcpublisher",
			"Messages":           msgList,
		}

		for _, comp := range d.CapabilityModel.Components {
			for _, telemetry := range comp.Telemetry {
				opcuaNodeId := fmt.Sprintf("nsu=%s;s=%s", d.getString(20), d.getString(20))
				telemetryName := telemetry.Name
				telemetryValue := d.getRandomValue(telemetry.Schema)
				payload[opcuaNodeId] = map[string]interface{}{
					"ServerTimestamp": time.Now(),
					"SourceTimestamp": time.Now(),
					"StatusCode":      nil,
					//"Name":            telemetryName,
					"Value": telemetryValue,
				}
				telemetryValues[telemetryName] = telemetryValue
				dataPointCount++
			}
		}

		body, err = json.Marshal(telemetryValues)
		if err != nil {
			return nil, err
		}
	} else {
		// typical device sending plain JSON payload confirming the DTDL model
		msg := make(map[string]interface{})
		for _, comp := range d.CapabilityModel.Components {
			for _, telemetry := range comp.Telemetry {
				msg[telemetry.Name] = d.getRandomValue(telemetry.Schema)
				dataPointCount++
			}
		}
		body, err = json.Marshal(msg)
		if err != nil {
			return nil, err
		}
	}
	correlationID, _ := uuid.GenerateUUID()
	messageID, _ := uuid.GenerateUUID()
	tm := telemetryMessage{
		body:               body,
		interfaceId:        "",
		connectionDeviceID: device.deviceID,
		connectionModuleID: "",
		contentEncoding:    "",
		contentType:        "Content-Type: application/json",
		correlationID:      correlationID,
		messageID:          messageID,
		creationTimeUtc:    creationTime, // distribute the messages in the batch evenly
		properties:         nil,
		dataPointCount:     dataPointCount,
	}
	telemetryMessages[0] = &tm

	return telemetryMessages, nil
}

// GenerateReportedProperties generate reported property update based on the device capability model.
func (d *DataGenerator) GenerateReportedProperties(device *device) (iotdevice.TwinState, error) {
	reportedProps := make(iotdevice.TwinState)
	for _, comp := range d.CapabilityModel.Components {
		for _, prop := range comp.Properties {
			if prop.Writable == false {
				reportedProps[prop.Name] = d.getRandomValue(prop.Schema)
			}
		}
	}
	return reportedProps, nil
}

// GenerateTwinUpdate creates a reported properties ACK based on the desired properties
func (d *DataGenerator) GenerateTwinUpdateAck(desiredTwin iotdevice.TwinState) iotdevice.TwinState {
	reportedTwin := make(iotdevice.TwinState)
	desiredVersion := desiredTwin.Version()
	for key, value := range desiredTwin {
		if key != "$version" {
			responseTwin := map[string]interface{}{
				"value": value,
				"ac":    200,
				"ad":    "completed",
				"av":    desiredVersion,
			}
			reportedTwin[key] = responseTwin

			values, ok := value.(map[string]interface{})
			if ok {
				_, ok := values["__t"]
				if ok {
					delete(values, "__t")

					componentTwin := map[string]interface{}{}
					componentTwin["__t"] = "c"
					for compKey, val := range values {
						componentTwin[compKey] = map[string]interface{}{
							"value": val,
							"ac":    200,
							"ad":    "completed",
							"av":    desiredVersion,
						}
					}

					reportedTwin[key] = componentTwin
				}
			}
		}
	}

	return reportedTwin
}

// GenerateTwinUpdate creates a reported properties ACK based on the desired properties
func (d *DataGenerator) GenerateC2DAck(c2dMsg *common.Message) *common.Message {
	var response common.Message

	return &response
}

func (d *DataGenerator) getRandomValue(schema string) interface{} {
	switch strings.ToLower(schema) {
	case "boolean":
		return d.getBool()
	case "date":
		return d.getDate()
	case "datetime":
		return d.getDateTime()
	case "double":
		return d.getDouble()
	case "duration":
		return d.getDuration()
	case "float":
		return d.getFloat()
	case "geopoint":
		return d.getGeopoint()
	case "integer":
		return d.getInt()
	case "long":
		return d.getLong()
	case "string":
		return d.getString(10)
	case "time":
		return d.getTime()
	}
	return ""
}

// getBool get a random boolean value.
func (d *DataGenerator) getBool() bool {
	return rand.Intn(100) < 50
}

// getDate gets the current date as a string.
func (d *DataGenerator) getDate() string {
	return time.Now().Format("2006-01-02")
}

// getDateTime gets current date time as a string.
func (d *DataGenerator) getDateTime() string {
	return time.Now().Format(time.RFC3339)
}

// getDouble gets a random double.
func (d *DataGenerator) getDouble() float64 {
	return 100 * rand.Float64()
}

// getDuration gets a random duration in ISO 8601 format.
func (d *DataGenerator) getDuration() string {
	// ISO 8601 format
	// P3Y6M4DT12H30M5S = three years, six months, four days, twelve hours, thirty minutes, and five seconds
	hr := rand.Int31n(12)
	min := rand.Int31n(60)
	sec := rand.Int31n(60)
	return fmt.Sprintf("P0Y0M0DT%dH%dM%dS", hr, min, sec)
}

// getFloat gets a random floating point number.
func (d *DataGenerator) getFloat() float32 {
	return 100 * rand.Float32()
}

// getInt gets a a geopoint along a predefined route in Redmond.
func (d *DataGenerator) getGeopoint() map[string]interface{} {
	d.nextGeoPoint = (d.nextGeoPoint + 1) % len(geopointRoute)
	return map[string]interface{}{
		"lat": geopointRoute[d.nextGeoPoint][0],
		"lon": geopointRoute[d.nextGeoPoint][1],
		"alt": 0,
	}
}

// getInt gets a random integer.
func (d *DataGenerator) getInt() int {
	return rand.Intn(100)
}

// getLong gets a random 64 bit integer.
func (d *DataGenerator) getLong() int64 {
	return rand.Int63n(1000)
}

// getString gets a random string.
func (d *DataGenerator) getString(length int) string {
	var charSet string = "abcdefghijklmnopqrstuvwxyzACBDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var val strings.Builder
	for i := 0; i < length; i++ {
		val.WriteString(string(charSet[rand.Intn(len(charSet))]))
	}

	return val.String()
}

//getTime gets the current time as string.
func (d *DataGenerator) getTime() string {
	return time.Now().Format(time.RFC3339)
}
