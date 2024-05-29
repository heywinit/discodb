package discodb

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Database struct {
	Token   string
	GuildId string
	Tables  []*Table
	Client  *discordgo.Session
}

type Column struct {
	Name   string
	Type   string
	Unique bool
}

type Table struct {
	Database *Database
	Name     string
	Columns  []Column
	Rows     []map[string]interface{}
}

func NewDatabase(token, guildId string) (*Database, error) {
	if token == "" || guildId == "" {
		return nil, errors.New("invalid Discord Bot Token or Guild ID")
	}

	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	err = client.Open()
	if err != nil {
		return nil, err
	}

	return &Database{
		Token:   token,
		GuildId: guildId,
		Tables:  make([]*Table, 0),
		Client:  client,
	}, nil
}

func (db *Database) CreateTable(name string, columns []Column) (*Table, error) {
	table := &Table{
		Name:    name,
		Columns: columns,
		Rows:    make([]map[string]interface{}, 0),
	}
	db.Tables = append(db.Tables, table)

	//TODO: implement a few channels under discodb category. those shall act as an internal storage. store the id of this channel in that channel
	_, err := db.Client.GuildChannelCreate(db.GuildId, name, discordgo.ChannelTypeGuildText)

	return table, err
}

func (db *Database) Close() error {
	return db.Client.Close()
}
