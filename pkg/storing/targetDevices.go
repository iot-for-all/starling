package storing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/iot-for-all/starling/pkg/models"
)

type targetDevices struct {
	store *store
}

// List lists all devices in a target from the store.
func (t *targetDevices) List(targetId string) ([]models.SimulationTargetDevice, error) {
	items := make([]models.SimulationTargetDevice, 0)
	prefix := []byte(fmt.Sprintf("targetDevices-%s-", targetId))
	err := t.store.list(prefix, func(k []byte, v []byte) error {
		var device models.SimulationTargetDevice
		err := json.Unmarshal(v, &device)
		if err != nil {
			return fmt.Errorf("failed to deserialize device %s: %w", k, err)
		}

		items = append(items, device)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// List lists all devices in a target from the store.
func (t *targetDevices) ListByTargetIdSimId(targetId string, simId string) ([]models.SimulationTargetDevice, error) {
	items := make([]models.SimulationTargetDevice, 0)
	prefix := []byte(fmt.Sprintf("targetDevices-%s-%s-%s", targetId, simId, targetId))
	err := t.store.list(prefix, func(k []byte, v []byte) error {
		var device models.SimulationTargetDevice
		err := json.Unmarshal(v, &device)
		if err != nil {
			return fmt.Errorf("failed to deserialize device %s: %w", k, err)
		}

		items = append(items, device)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Get gets a specific target's devices.
func (t *targetDevices) Get(targetId string, deviceId string) (*models.SimulationTargetDevice, error) {
	var item models.SimulationTargetDevice

	err := t.store.get([]byte(fmt.Sprintf("targetDevices-%s-%s", targetId, deviceId)), &item)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Set create or updates the device in a target.
func (t *targetDevices) Set(item *models.SimulationTargetDevice) error {
	return t.store.set([]byte(fmt.Sprintf("targetDevices-%s-%s", item.TargetID, item.DeviceID)), item)
}

// Delete deletes device in a target.
func (t *targetDevices) Delete(targetId string, deviceId string) error {
	err := t.store.delete([]byte(fmt.Sprintf("targetDevices-%s-%s", targetId, deviceId)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

// DeleteAll deletes all devices in a target.
func (t *targetDevices) DeleteAll(targetId string) error {
	err := t.store.delete([]byte(fmt.Sprintf("targetDevices-%s-", targetId)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
