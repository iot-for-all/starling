package storing

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/iot-for-all/starling/pkg/models"
)

// deviceConfigs represents the device configurations of a simulation
type deviceConfigs struct {
	store *store
}

// Get get a device config for the given simulation and config id
func (s *deviceConfigs) Get(simulationID string, configID string) (*models.SimulationDeviceConfig, error) {
	var item models.SimulationDeviceConfig
	err := s.store.get([]byte(fmt.Sprintf("deviceConfig-%s-%s", simulationID, configID)), &item)

	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// List lists all existing device configs
func (s *deviceConfigs) List(simulationID string) ([]*models.SimulationDeviceConfig, error) {
	items := make([]*models.SimulationDeviceConfig, 0)
	prefix := []byte(fmt.Sprintf("deviceConfig-%s-", simulationID))
	err := s.store.list(prefix, func(k []byte, v []byte) error {
		var cfg models.SimulationDeviceConfig
		err := json.Unmarshal(v, &cfg)
		if err != nil {
			return fmt.Errorf("failed to deserialize simulation device %s: %w", k, err)
		}

		items = append(items, &cfg)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Set create or updates the device config
func (s *deviceConfigs) Set(simulationID string, config *models.SimulationDeviceConfig) error {
	return s.store.set([]byte(fmt.Sprintf("deviceConfig-%s-%s", simulationID, config.ID)), config)
}

// Delete deletes an existing device config
func (s *deviceConfigs) Delete(simulationID string, configID string) error {
	err := s.store.delete([]byte(fmt.Sprintf("deviceConfig-%s-%s", simulationID, configID)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	return err
}
