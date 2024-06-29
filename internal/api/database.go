package api

import (
	"DiscoDB/internal/models"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// DBClient represents a client for the Discord-based database.
type DBClient struct {
	Session *discordgo.Session
}

// NewDBClient creates a new DBClient instance.
func NewDBClient(token string) (*DBClient, error) {
	intents := discordgo.IntentsGuilds |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuildMessageTyping |
		discordgo.IntentsGuildVoiceStates |
		discordgo.IntentsGuildPresences

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = intents

	client := &DBClient{
		Session: session,
	}

	// Open the Discord session
	err = session.Open()
	if err != nil {
		return nil, err
	}

	return client, nil
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
	var metadataChannelID string
	channels, _ := client.Session.GuildChannels(guild.ID)
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory && channel.Name == "internal" {
			internalCategoryID = channel.ID
		}
		if channel.Type == discordgo.ChannelTypeGuildText && channel.Name == "tables" && channel.ParentID == internalCategoryID {
			metadataChannelID = channel.ID
			//get metadata for tables
			messages, _ := client.Session.ChannelMessages(metadataChannelID, 1, "", "", "")
			for _, message := range messages {
				var metadata map[string]interface{}
				err := json.Unmarshal([]byte(message.Content), &metadata)
				if err != nil {
					return nil, err
				}

				//get tables in metadata
				tables := metadata["tables"].(map[string]interface{})
				for _, table := range tables {
					tableMap, ok := table.(map[string]interface{})
					if !ok {
						continue
					}

					channelId := tableMap["id"].(string)
					name := tableMap["name"]
					schema := tableMap["schema"]
					fmt.Println(channelId, name, schema)
				}
			}
		}
	}

	db := &models.Database{
		ID:                 guild.ID,
		Name:               guild.Name,
		InternalCategoryID: internalCategoryID,
		MetadataChannelID:  metadataChannelID,
	}

	return db, nil
}

// DeleteDatabase deletes an existing Discord server representing a database.
func (client *DBClient) DeleteDatabase(databaseID string) error {
	err := client.Session.GuildDelete(databaseID)
	if err != nil {
		return err
	}
	return nil
}

func (client *DBClient) Close() error {
	if client.Session != nil {
		err := client.Session.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
