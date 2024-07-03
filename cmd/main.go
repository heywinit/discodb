package main

import (
	"DiscoDB/internal/api"
	"DiscoDB/internal/models"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN must be set")
	}

	client, err := api.NewDBClient(token)
	if err != nil {
		log.Fatal("Error creating Discord client:", err)
	}
	defer func(client *api.DBClient) {
		err := client.Close()
		if err != nil {
			log.Fatal("Error closing Database client:", err)
		}
	}(client)

	//// Get all guilds the bot is connected to
	//guilds, err := client.Session.UserGuilds(100, "", "", false)
	//if err != nil {
	//	log.Fatal("Error fetching guilds:", err)
	//}
	//
	//log.Printf("Connected to %d guilds\n", len(guilds))
	////this is for testing purposes where we sometimes create too many databases
	//for _, guild := range guilds {
	//	if guild.Owner {
	//		err := client.Session.GuildDelete(guild.ID)
	//		if err != nil {
	//			log.Fatal("Error deleting guild:", err)
	//		}
	//	}
	//}

	// Example usage: Create a new database
	database, err := client.LoadDatabase(os.Getenv("DISCORD_GUILD_ID"))
	//database, err := client.CreateDatabase("TestDatabase2")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}

	schema := map[string]string{
		"id":    "int",
		"name":  "string",
		"email": "string",
	}

	// Example usage: Create a new table (text channel)
	table, err := client.CreateTable(*database, "users", schema)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	// Example usage: Create a new record in the table
	newRecord := models.Record{
		ID:     "1",
		Fields: map[string]interface{}{"id": 1, "name": "Alice"},
	}

	err = client.CreateRecord(table.ID, newRecord)
	if err != nil {
		log.Fatal("Error creating record:", err)
	}

	// Example usage: Read all records from the table
	records, err := client.ReadRecords(table.ID)
	if err != nil {
		log.Fatal("Error reading records:", err)
	}

	// Create an invite link for the table (text channel)
	invite, err := client.Session.ChannelInviteCreate(table.ID, discordgo.Invite{
		MaxAge:  0, // Invite link does not expire
		MaxUses: 0, // Unlimited uses
	})
	if err != nil {
		log.Fatal("Error creating invite link:", err)
	}

	log.Printf("Invite link: https://discord.gg/%s\n", invite.Code)
	log.Printf("Records read from table %s: %v\n", table.Name, records)
}
