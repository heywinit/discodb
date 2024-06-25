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
	// get id store
	idStore, err := LoadSavedDatabase("store.didb")
	if err != nil {
		return nil, err
	}

	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	err = client.Open()
	if err != nil {
		return nil, err
	}

	// check if the categories and channels exist
	x := []string{"tables", "discodb"}
	for i := 0; i < len(x); i++ {
		if _, ok := idStore.Categories[x[i]]; !ok {
			tableCategory, err := client.GuildChannelCreate(guildId, x[i], discordgo.ChannelTypeGuildCategory)
			if err != nil {
				return nil, err
			}
			idStore.Categories[x[i]] = tableCategory.ID
		}
	}

	// check if channels exist
	tables := make([]*Table, 0)
	x = []string{"tables"}
	for i := 0; i < len(x); i++ {
		if _, ok := idStore.Channels[x[i]]; !ok {
			tablesChannel, err := client.GuildChannelCreate(guildId, x[i], discordgo.ChannelTypeGuildText)
			if err != nil {
				return nil, err
			}
			_, err = client.ChannelEditComplex(tablesChannel.ID, &discordgo.ChannelEdit{
				ParentID: idStore.Categories["discodb"],
			})
			if err != nil {
				return nil, err
			}
			idStore.Channels[x[i]] = tablesChannel.ID
		} else {
			// get message from this tables channel
			tablesChannel, err := client.Channel(idStore.Channels[x[i]])
			if err != nil {
				return nil, err
			}
			messages, err := client.ChannelMessages(tablesChannel.ID, 1, "", "", "")
			if err != nil {
				return nil, err
			}
			// if no messages exist
			if len(messages) == 0 {
				tables := make([]string, 0)
				tablesJson, _ := json.Marshal(tables)
				_, err = client.ChannelMessageSend(tablesChannel.ID, string(tablesJson))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// deserialize the tables
	for _, channel := range idStore.Channels {
		tablesChannel, err := client.Channel(channel)
		if err != nil {
			return nil, err
		}
		messages, err := client.ChannelMessages(tablesChannel.ID, 1, "", "", "")
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			var deserializedTables []string
			err = json.Unmarshal([]byte(messages[0].Content), &deserializedTables)
			if err != nil {
				return nil, err
			}
			for i := 0; i < len(deserializedTables); i++ {
				var table Table
				err = json.Unmarshal([]byte(deserializedTables[i]), &table)
				if err != nil {
					return nil, err
				}
				// append the table to the database
				tables = append(tables, &table)
			}
		}
	}

	return &Database{
		Token:   token,
		GuildId: guildId,
		Tables:  tables,
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
