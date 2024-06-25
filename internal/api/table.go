package api

import (
	"DiscoDB/internal/models"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
)

// CreateTable creates a new text channel to represent a table within a specific database.
func (client *DBClient) CreateTable(database models.Database, tableName string, schema map[string]string) (*models.Table, error) {
	internalCategoryID := database.InternalCategoryID
	if internalCategoryID == "" {
		return nil, errors.New("internal category not found")
	}

	// Get the tables metadata channel
	tablesMetadataChannelID := database.MetadataChannelID
	if tablesMetadataChannelID == "" {
		return nil, errors.New("tables metadata channel not found")
	}

	// Create the new table (text channel)
	channel, err := client.Session.GuildChannelCreate(database.ID, tableName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	table := &models.Table{
		ID:     channel.ID,
		Name:   channel.Name,
		Schema: schema,
	}

	// Update metadata
	messages, err := client.Session.ChannelMessages(tablesMetadataChannelID, 1, "", "", "")
	if err != nil || len(messages) == 0 {
		return nil, errors.New("failed to retrieve metadata message")
	}

	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(messages[0].Content), &metadata)
	if err != nil {
		return nil, err
	}

	tables := metadata["tables"].(map[string]interface{})
	tables[tableName] = table
	metadata["tables"] = tables

	updatedMetadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	_, err = client.Session.ChannelMessageEdit(tablesMetadataChannelID, messages[0].ID, string(updatedMetadataJSON))
	if err != nil {
		return nil, err
	}

	return table, nil
}

// DeleteTable deletes an existing table (text channel) from a database.
func (client *DBClient) DeleteTable(databaseID string, tableID string) error {
	_, err := client.Session.ChannelDelete(tableID)
	if err != nil {
		return err
	}
	return nil
}
