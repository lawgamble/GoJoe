package main

import (
	"fmt"
	"goJoe/internal/facade"
	"goJoe/internal/facade/vouch"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var i int

func main() {
	err := godotenv.Load("./local.env")
	if err != nil {
		fmt.Println(err)
	}

	bot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN")) //figure out how to get the token from .env
	if err != nil {
		panic(err)
	}
	// register events
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	err = bot.Open() // creates web socket to Discord
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	for {
		if i == 1 {
			time.Sleep(time.Second * 3)
			return
		}
	} // this for loop keeps the bot running until i == 1
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("WE ARE LIVE!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == os.Getenv("BOT_ID") {
		return
	}

	command := strings.Fields(m.Content)[0]

	switch command {
	case "!gostatus":
		{
			s.ChannelMessageSend(m.ChannelID, "I'm alive")
		}
	case "!enditall":
		{
			i = 1
			s.ChannelMessageSend(m.ChannelID, "Goodbye!")
			return
		}
	case "!register":
		{
			res := facade.Register(m)
			s.ChannelMessageSend(m.ChannelID, res)
		}
	case "!vouch":
		{
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			vouch.Vouch(s, m)
		}
	case "!unvouch":
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		vouch.UnVouch(s, m)
	}
}
