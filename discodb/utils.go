package discodb

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"os"
)

func LoadSavedDatabase(filename string) (*IDStore, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var store IDStore
	err = json.Unmarshal(file, &store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (db *Database) SaveDatabase(filename string, store IDStore) error {
	storeJson, _ := json.Marshal(store)
	return os.WriteFile(filename, storeJson, 0644)
}

// serialize converts the database to a binary representation
func (db *Database) serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(db); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
