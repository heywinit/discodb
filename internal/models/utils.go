package models

import (
	"encoding/json"
)

// MarshalRecord converts a Record to a JSON string.
func MarshalRecord(record *Record) (string, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalRecord converts a JSON string to a Record.
func UnmarshalRecord(data string) (*Record, error) {
	var record Record
	err := json.Unmarshal([]byte(data), &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// MarshalTable converts a Table to a JSON string.
func MarshalTable(table *Table) (string, error) {
	data, err := json.Marshal(table)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalTable converts a JSON string to a Table.
func UnmarshalTable(data string) (*Table, error) {
	var table Table
	err := json.Unmarshal([]byte(data), &table)
	if err != nil {
		return nil, err
	}
	return &table, nil
}
