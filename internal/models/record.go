package models

// Record represents a record in a table.
type Record struct {
	ID     string                 `json:"id"`     // Message ID in Discord
	Fields map[string]interface{} `json:"fields"` // Field name -> Value
}
