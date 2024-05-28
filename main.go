package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"log"
)

func main() {
	env, _ := godotenv.Read(".env")
	token := env["DISCORD_TOKEN"]
	if token == "" {
		log.Fatal("No token provided")
		return
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session: ", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//ignore bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			return
		}
	}
}
