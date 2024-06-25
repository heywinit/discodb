package main

import (
	"DiscoDB/internal/api"
	"log"
	"os"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN must be set")
	}

	client, err := api.NewDBClient(token)
	if err != nil {
		log.Fatal("Error creating Discord client:", err)
	}

	// Example usage
	database, err := client.CreateDatabase("MyDatabase")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}

	schema := map[string]string{
		"id":   "int",
		"name": "string",
	}

	table, err := client.CreateTable(database.ID, "Users", schema)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	log.Printf("Database created: %v\n", database)
	log.Printf("Table created: %v\n", table)
}
