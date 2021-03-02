package storing

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/iot-for-all/starling/pkg/models"
)

type deviceModels struct {
	store *store
}

// Get gets the device model for the given id.
func (m *deviceModels) Get(id string) (*models.DeviceModel, error) {
	var model models.DeviceModel
	err := m.store.get([]byte(fmt.Sprintf("deviceModel-%s", id)), &model)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &model, nil
}

// List lists all existing device models.
func (m *deviceModels) List() ([]models.DeviceModel, error) {
	items := make([]models.DeviceModel, 0)
	prefix := []byte("deviceModel-")
	err := m.store.list(prefix, func(k []byte, v []byte) error {
		var model models.DeviceModel
		err := json.Unmarshal(v, &model)
		if err != nil {
			return fmt.Errorf("failed to deserialize simulation: %w", err)
		}

		items = append(items, model)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Set creates or updates a device model.
func (m *deviceModels) Set(item *models.DeviceModel) error {
	return m.store.set([]byte(fmt.Sprintf("deviceModel-%s", item.ID)), item)
}

// Delete deletes an existing device model.
func (m *deviceModels) Delete(id string) error {
	err := m.store.delete([]byte(fmt.Sprintf("deviceModel-%s", id)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
