package main

import (
	"DiscoDB/internal/api"
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

	// Example usage
	database, err := client.LoadDatabase(os.Getenv("DISCORD_GUILD_ID"))
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

	// Create an invite link for the table (text channel)
	invite, err := client.Session.ChannelInviteCreate(table.ID, discordgo.Invite{
		MaxAge:  0, // Invite link does not expire
		MaxUses: 0, // Unlimited uses
	})
	if err != nil {
		log.Fatal("Error creating invite link:", err)
	}

	log.Printf("Database created: %v\n", database)
	log.Printf("Table created: %v\n", table)
	log.Printf("Invite link: https://discord.gg/%s\n", invite.Code)
}
