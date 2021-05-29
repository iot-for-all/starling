package storing

import (
	"encoding/json"
	"fmt"
	"github.com/iot-for-all/starling/pkg/config"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"
)

var (
	db            *badger.DB     // application database
	DeviceModels  *deviceModels  // DeviceModels store
	Simulations   *simulations   // Simulations store
	DeviceConfigs *deviceConfigs // DeviceConfigs store
	Targets       *targets       // Targets store
	TargetModels  *targetModels  // TargetModels store
	TargetDevices *targetDevices // TargetDevices store
)

type store struct {
	db *badger.DB
}

// Open initializes and opens the database
func Open(cfg *config.StoreConfig) error {
	// TODO: Open with correct badger options
	dbFile := fmt.Sprintf("%s/.db", cfg.DataDirectory)
	opts := badger.DefaultOptions(dbFile)
	//b, _ := json.MarshalIndent(opts, "", "  ")
	//fmt.Printf("%s\n", b)
	d, err := badger.Open(opts)
	if err != nil {
		return err
	}

	db = d
	store := store{db: db}

	DeviceModels = &deviceModels{store: &store}
	Simulations = &simulations{store: &store}
	DeviceConfigs = &deviceConfigs{store: &store}
	Targets = &targets{store: &store}
	TargetModels = &targetModels{store: &store}
	TargetDevices = &targetDevices{store: &store}

	log.Info().Msgf("initialized database from %s", dbFile)
	return nil
}

// Close closes the open database handle
func Close() error {
	if db != nil {
		return db.Close()
	}

	return nil
}

// get gets the value of a key
func (s *store) get(key []byte, target interface{}) error {
	return s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(v []byte) error {
			if err := json.Unmarshal(v, target); err != nil {
				return fmt.Errorf("error de-serializing value for %s from store", key)
			}

			return nil
		})

		return nil
	})
}

// list lists all existing rows that match the key prefix
func (s *store) list(prefix []byte, handler func(key []byte, val []byte) error) error {
	return s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10

		it := txn.NewIterator(opts)

		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				return handler(item.Key(), v)
			})

			if err != nil {
				return fmt.Errorf("failure reading items for key %s, %w", prefix, err)
			}
		}

		return nil
	})
}

// set create or updates the specified value
func (s *store) set(key []byte, target interface{}) error {
	return s.db.Update(func(txn *badger.Txn) error {
		val, err := json.Marshal(target)
		if err != nil {
			return fmt.Errorf("failed to serialize %s: %w", key, err)
		}

		if err = txn.Set(key, val); err != nil {
			return fmt.Errorf("failed to save %s: %w", key, err)
		}

		return nil
	})
}

// delete deletes the value for the specified key
func (s *store) delete(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(key)
		if err != nil {
			return fmt.Errorf("failed to delete key %s: %w", key, err)
		}

		return nil
	})
}
