package api

import (
	"DiscoDB/internal/models"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)

// DBClient represents a client for the Discord-based database.
type DBClient struct {
	Session *discordgo.Session
}

// NewDBClient creates a new DBClient instance.
func NewDBClient(token string) (*DBClient, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &DBClient{Session: session}, nil
}

// CreateDatabase creates a new Discord server to represent a database.
func (client *DBClient) CreateDatabase(name string) (*models.Database, error) {
	guild, err := client.Session.GuildCreate(name)
	if err != nil {
		return nil, err
	}

	// Create internal category
	category, err := client.Session.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     "internal",
		Type:     discordgo.ChannelTypeGuildCategory,
		Position: 0,
	})
	if err != nil {
		return nil, err
	}

	// Create metadata channel for table info
	tableMetadataChannel, err := client.Session.GuildChannelCreateComplex(guild.ID, discordgo.GuildChannelCreateData{
		Name:     "tables",
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		Topic:    "Metadata for tables",
	})
	if err != nil {
		return nil, err
	}

	// Initialize metadata
	initialMetadata := map[string]interface{}{
		"tables": make(map[string]interface{}),
	}
	metadataJSON, err := json.Marshal(initialMetadata)
	if err != nil {
		return nil, err
	}

	_, err = client.Session.ChannelMessageSend(tableMetadataChannel.ID, string(metadataJSON))
	if err != nil {
		return nil, err
	}

	return &models.Database{
		ID:                 guild.ID,
		Name:               guild.Name,
		InternalCategoryID: category.ID,
		MetadataChannelID:  tableMetadataChannel.ID,
	}, nil
}

// LoadDatabase loads an existing Discord server (database) by its ID.
func (client *DBClient) LoadDatabase(databaseID string) (*models.Database, error) {
	guild, err := client.Session.Guild(databaseID)
	if err != nil {
		return nil, err
	}

	var internalCategoryID string
	for _, channel := range guild.Channels {
		if channel.Name == "internal" && channel.Type == discordgo.ChannelTypeGuildCategory {
			internalCategoryID = channel.ID
			break
		}
	}

	return &models.Database{
		ID:                 guild.ID,
		Name:               guild.Name,
		InternalCategoryID: internalCategoryID,
	}, nil
}

// DeleteDatabase deletes an existing Discord server representing a database.
func (client *DBClient) DeleteDatabase(databaseID string) error {
	err := client.Session.GuildDelete(databaseID)
	if err != nil {
		return err
	}
	return nil
}
