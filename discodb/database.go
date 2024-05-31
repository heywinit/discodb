package discodb

import (
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
)

type IDStore struct {
	Categories map[string]string `json:"categories"`
	Channels   map[string]string `json:"channels"`
}

type Database struct {
	Token   string
	GuildId string
	Tables  []*Table
	Client  *discordgo.Session
	IdStore IDStore
}

type Column struct {
	Name   string
	Type   string
	Unique bool
}

// Table TODO: now that im thinking of it, we might use something like mongo does with bson. so you can give a struct, and it serializes that into something and stores that. but that would add overhead. anyways, its 1am i should go sleep
type Table struct {
	Database *Database                `json:"database"`
	Name     string                   `json:"name"`
	Columns  []Column                 `json:"columns"`
	Rows     []map[string]interface{} `json:"rows"`
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

	idStore := IDStore{}
	categories := make(map[string]string)

	tableCategory, err := client.GuildChannelCreate(guildId, "tables", discordgo.ChannelTypeGuildCategory)
	if err != nil {
		return nil, err
	}
	categories["tables"] = tableCategory.ID

	discoDBCategory, err := client.GuildChannelCreate(guildId, "discodb", discordgo.ChannelTypeGuildCategory)
	if err != nil {
		return nil, err
	}
	categories["discodb"] = discoDBCategory.ID

	channels := make(map[string]string)
	tablesChannel, err := client.GuildChannelCreate(guildId, "tables", discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}
	_, err = client.ChannelEditComplex(tablesChannel.ID, &discordgo.ChannelEdit{
		ParentID: discoDBCategory.ID,
	})
	if err != nil {
		return nil, err
	}
	channels["tables"] = tablesChannel.ID

	idStore.Categories = categories
	idStore.Channels = channels

	return &Database{
		Token:   token,
		GuildId: guildId,
		Tables:  make([]*Table, 0),
		Client:  client,
		IdStore: idStore,
	}, nil
}

func LoadDatabase(token, guildId string) (*Database, error) {
	//get id store
	idStore, err := LoadSavedDatabase("store.didb")
	if err != nil {
		panic(err)
	}

	client, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	err = client.Open()
	if err != nil {
		panic(err)
	}

	return &Database{
		Token:   token,
		GuildId: guildId,
		Tables:  make([]*Table, 0),
		Client:  client,
		IdStore: *idStore,
	}, nil
}

func (db *Database) CreateTable(name string, columns []Column) (*Table, error) {
	table := &Table{
		Name:    name,
		Columns: columns,
		Rows:    make([]map[string]interface{}, 0),
	}
	db.Tables = append(db.Tables, table)

	createdTable, err := db.Client.GuildChannelCreate(db.GuildId, name, discordgo.ChannelTypeGuildText)
	if err != nil {
		return nil, err
	}
	_, err = db.Client.ChannelEditComplex(createdTable.ID, &discordgo.ChannelEdit{
		ParentID: db.IdStore.Categories["tables"],
	})
	if err != nil {
		return nil, err
	}

	//serialize the created table into json
	tableJson, _ := json.Marshal(table)

	//get tables channel
	tablesChannel, err := db.Client.Channel(db.IdStore.Channels["tables"])
	if err != nil {
		return nil, err
	}
	messages, err := db.Client.ChannelMessages(tablesChannel.ID, 1, "", "", "")
	if err != nil {
		return nil, err
	}

	//if no other tables exist
	if len(messages) == 0 {
		tables := []string{string(tableJson)}
		tablesJson, _ := json.Marshal(tables)
		_, err = db.Client.ChannelMessageSend(tablesChannel.ID, string(tablesJson))
		if err != nil {
			return nil, err
		}
	}

	//if other tables exist
	if len(messages) > 0 {
		var tables []string
		err = json.Unmarshal([]byte(messages[0].Content), &tables)
		if err != nil {
			return nil, err
		}
		tables = append(tables, string(tableJson))
		tablesJson, _ := json.Marshal(tables)
		_, err = db.Client.ChannelMessageEdit(tablesChannel.ID, messages[0].ID, string(tablesJson))
		if err != nil {
			return nil, err
		}
	}

	db.Tables = append(db.Tables, table)
	return table, err
}

func (db *Database) Close() error {
	//store the IdStore into a file
	err := db.SaveDatabase("store.didb", db.IdStore)
	if err != nil {
		return err
	}

	err = db.Client.Close()
	if err != nil {
		return err
	}

	return nil
}
