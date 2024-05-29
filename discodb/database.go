package discodb

import "errors"

type Database struct {
	Token   string
	GuildId string
	Tables  []*Table
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
		return nil, errors.New("token and guildId cannot be empty")
	}

	return &Database{
		Token:   token,
		GuildId: guildId,
		Tables:  make([]*Table, 0),
	}, nil
}

func (db *Database) CreateTable(name string, columns []Column) (*Table, error) {
	table := &Table{
		Name:    name,
		Columns: columns,
		Rows:    make([]map[string]interface{}, 0),
	}
	db.Tables = append(db.Tables, table)
	return table, nil
}
