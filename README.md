# DiscoDB
A database that stores data in Discord. Created using GoLang & [discordgo](https://github.com/bwmarrin/discordgo).

[![AGPL License](https://img.shields.io/github/license/vineshgh/DiscoDB)](https://img.shields.io/github/license/vineshgh/DiscoDB/LICENSE)


## Usage/Examples
```golang
db, err := discodb.NewDatabase(token, guildId)	

table, err = db.CreateTable("users", []discodb.Column{
    {Name: "id", Type: "int", Unique: true},
    {Name: "name", Type: "string"},
    {Name: "age", Type: "int"},
})

db.Close()
```

## Features
- [x] Unlimited data storage
- [ ] SDK/API for easy comms
- [ ] Data encryption


## FAQ
#### Whats the ETA?
I am learning Go as I make this. This is a way for me to understand how real life things are done in GoLang. So I am not sure about the ETA

