package serving

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/iot-for-all/starling/pkg/models"
	"io"
	"net/http"
	"strings"
	"time"
)

type DeviceTemplateDownloader struct {
	target *models.SimulationTarget
}

func NewDeviceTemplateDownloader(target *models.SimulationTarget) *DeviceTemplateDownloader {
	return &DeviceTemplateDownloader{
		target: target,
	}
}

func (d *DeviceTemplateDownloader) DownloadModels() ([]*models.DeviceModel, error) {
	var deviceModels []*models.DeviceModel
	ctx := context.Background()
	path := fmt.Sprintf("https://%s/api/deviceTemplates?api-version=1.0", d.target.AppUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return deviceModels, err
	}

	req.Header.Add("Authorization", d.target.AppToken)
	client := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return deviceModels, err

	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return deviceModels, err
	}

	var dtr models.DeviceTemplateResponse
	err = json.Unmarshal(body, &dtr)
	if err != nil {
		return deviceModels, err
	}

	for _, deviceTemplate := range dtr.Value {
		capabilityModel, ok := deviceTemplate["capabilityModel"].(map[string]interface{})
		if !ok {
			return deviceModels, err
		}
		name, ok := deviceTemplate["displayName"].(string)
		if !ok {
			return deviceModels, fmt.Errorf("could not find displayName in capability model")
		}
		name = d.scrubModelName(name)
		model := models.DeviceModel{
			ID:              name,
			Name:            name,
			CapabilityModel: []map[string]interface{}{capabilityModel},
		}
		deviceModels = append(deviceModels, &model)
	}

	return deviceModels, nil
}

func (d *DeviceTemplateDownloader) scrubModelName(name string) string {
	const alphaNumeric = "abcdefghijklmnopqrstuvwxyz0123456789"
	buffer := bytes.Buffer{}
	name = strings.ToLower(name)
	for _, c := range name {
		if strings.Contains(alphaNumeric, string(c)) {
			buffer.WriteString(string(c))
		}
	}
	return buffer.String()
}
