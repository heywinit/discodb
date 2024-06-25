package discodb

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

// DropTable
func (db *Database) DropTable(name string) error {
	for i, table := range db.Tables {
		if table.Name == name {
			db.Tables = append(db.Tables[:i], db.Tables[i+1:]...)
			return nil
		}
	}
	return nil
}
