package api

import (
	"DiscoDB/internal/models"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)

// CreateTable creates a new text channel to represent a table within a specific database.
func (client *DBClient) CreateTable(databaseID string, tableName string, schema map[string]string) (*models.Table, error) {
	channel, err := client.Session.GuildChannelCreate(databaseID, tableName, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}

	table := &models.Table{
		ID:     channel.ID,
		Name:   channel.Name,
		Schema: schema,
	}

	// Store the schema in the channel topic or a pinned message
	schemaData, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	_, err = client.Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Topic: string(schemaData),
	})
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
