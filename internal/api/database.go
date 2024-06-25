package api

import (
	"DiscoDB/internal/models"
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
	return &models.Database{
		ID:   guild.ID,
		Name: guild.Name,
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
