package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"log"
	"os"
	"os/signal"
)

var token string
var s *discordgo.Session

func init() {
	env, _ := godotenv.Read(".env")
	token = env["DISCORD_TOKEN"]
	var err error
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	s.Identify.Intents = discordgo.IntentsGuildMessages

	err := s.Open()
	if err != nil {
		log.Fatal("Error opening connection: ", err)
		return
	}

	defer s.Close()

	log.Println("Bot is now running.")

	// Keep the bot running
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
