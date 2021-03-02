package storing

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/iot-for-all/starling/pkg/models"
)

type targetModels struct {
	store *store
}

// Get gets a specific target's models.
func (t *targetModels) Get(targetId string) (*models.SimulationTargetModels, error) {
	var item models.SimulationTargetModels

	err := t.store.get([]byte(fmt.Sprintf("targetModels-%s", targetId)), &item)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// Set creates or updates the models configured for a target.
func (t *targetModels) Set(item *models.SimulationTargetModels) error {
	return t.store.set([]byte(fmt.Sprintf("targetModels-%s", item.TargetID)), item)
}

// Delete deletes models configured for a target.
func (t *targetModels) Delete(targetId string) error {
	err := t.store.delete([]byte(fmt.Sprintf("targetModels-%s", targetId)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
