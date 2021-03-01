package storing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/reddyduggempudi/starling/pkg/models"
)

type targets struct {
	store *store
}

// Get gets a specific target from the store by its id.
func (t *targets) Get(id string) (*models.SimulationTarget, error) {
	var item models.SimulationTarget
	err := t.store.get([]byte(fmt.Sprintf("target-%s", id)), &item)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &item, nil
}

// List lists all targets in the store.
func (t *targets) List() ([]models.SimulationTarget, error) {
	items := make([]models.SimulationTarget, 0)
	prefix := []byte("target-")
	err := t.store.list(prefix, func(k []byte, v []byte) error {
		var target models.SimulationTarget
		err := json.Unmarshal(v, &target)
		if err != nil {
			return fmt.Errorf("failed to deserialize target %s: %w", k, err)
		}

		items = append(items, target)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Set creates or updates a target
func (t *targets) Set(item *models.SimulationTarget) error {
	return t.store.set([]byte(fmt.Sprintf("target-%s", item.ID)), item)
}

// Delete deletes an existing target
func (t *targets) Delete(id string) error {
	err := t.store.delete([]byte(fmt.Sprintf("target-%s", id)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
