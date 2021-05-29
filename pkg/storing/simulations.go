package storing

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/iot-for-all/starling/pkg/models"
	"github.com/rs/zerolog/log"
)

type simulations struct {
	store *store
}

// Get gets a simulation by its id.
func (s *simulations) Get(id string) (*models.Simulation, error) {
	var item models.Simulation
	err := s.store.get([]byte(fmt.Sprintf("simulation-%s", id)), &item)
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to Get simulation")
		return nil, err
	}

	return &item, nil
}

// List lists all existing simulations
func (s *simulations) List() ([]models.Simulation, error) {
	items := make([]models.Simulation, 0)
	prefix := []byte("simulation-")
	err := s.store.list(prefix, func(k []byte, v []byte) error {
		var sim models.Simulation
		err := json.Unmarshal(v, &sim)
		if err != nil {
			return fmt.Errorf("failed to deserialize simulation %s, %s", k, err.Error())
		}

		items = append(items, sim)
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to list simulations")
		return nil, err
	}

	return items, nil
}

// Set creates or updates a simulation
func (s *simulations) Set(item *models.Simulation) error {
	return s.store.set([]byte(fmt.Sprintf("simulation-%s", item.ID)), item)
}

// Delete deletes an existing simulation
func (s *simulations) Delete(id string) error {
	err := s.store.delete([]byte(fmt.Sprintf("simulation-%s", id)))
	if err != nil && errors.Is(err, badger.ErrKeyNotFound) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
