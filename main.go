package main

import (
	"DiscoDB/discodb"
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("TOKEN")
	guildId := os.Getenv("GUILD_ID")

	fmt.Printf("   Token: %s\n", token)
	fmt.Printf("Guild ID: %s\n", guildId)

	db, err := discodb.NewDatabase(token, guildId)
	if err != nil {
		panic(err)
	}

	_, err = db.CreateTable("users", []discodb.Column{
		{Name: "id", Type: "int", Unique: true},
		{Name: "name", Type: "string"},
		{Name: "age", Type: "int"},
	})

	db.Close()
}
