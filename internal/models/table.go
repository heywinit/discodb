package models

// Table represents a text channel acting as a table within a database.
type Table struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Schema map[string]string `json:"schema"` // Field name -> Data type
}
