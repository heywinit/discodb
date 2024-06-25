package models

// Database represents a Discord server acting as a database.
type Database struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	InternalCategoryID string `json:"internal_category_id"`
	MetadataChannelID  string `json:"metadata_channel_id"`
}
