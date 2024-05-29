package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"log"
	"os"
	"os/signal"
)

var token string
var guildId string
var channelId string
var configMsgId string
var s *discordgo.Session

func init() {
	env, _ := godotenv.Read(".env")
	token = env["DISCORD_TOKEN"]
	guildId = env["DISCORD_GUILD_ID"]
	channelId = env["DISCORD_CHANNEL_ID"]
	configMsgId = env["DISCORD_CONFIG_MSG_ID"]
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

	loadDatabase()

	// Keep the bot running
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	shutdown()
}

func loadDatabase() {
	log.Printf("Loading database. Guild: %v\n", guildId)
	_, err := s.Guild(guildId)
	if err != nil {
		log.Fatalf("Error loading guild: %v", err)
		return
	}

	if channelId == "" {
		//create category discodb
		cat, err := s.GuildChannelCreate(guildId, "discodb", discordgo.ChannelTypeGuildCategory)
		//create channel dbconfig
		ch, err := s.GuildChannelCreate(guildId, "dbconfig", discordgo.ChannelTypeGuildText)
		//move channel to category
		_, err = s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
			ParentID: cat.ID,
		})
		if err != nil {
			log.Fatalf("Error creating channel: %v", err)
			return
		}

		channelId = ch.ID
		msg, err := s.ChannelMessageSend(channelId, "Config message")
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
			return
		}

		channelId = ch.ID
		configMsgId = msg.ID
	} else {
		_, err := s.Channel(channelId)
		if err != nil {
			log.Fatalf("Error loading channel: %v", err)
			return
		}
	}
	//load config message
	msg, err := s.ChannelMessage(channelId, configMsgId)
	if err != nil {
		log.Fatalf("Error loading message: %v", err)
		return
	}
	fmt.Printf("Message: %v\n", msg.Content)
}

func shutdown() {
	env, _ := godotenv.Read(".env")
	env["DISCORD_TOKEN"] = token
	env["DISCORD_GUILD_ID"] = guildId
	env["DISCORD_CHANNEL_ID"] = channelId
	env["DISCORD_CONFIG_MSG_ID"] = configMsgId
	godotenv.Write(env, ".env")
}
