package main

import (
	"fmt"
	"goJoe/internal/facade"
	"goJoe/internal/facade/repeat"
	"goJoe/internal/facade/status"
	"goJoe/internal/facade/tripleCrown"
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

	bot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

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

func ready(s *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("WE ARE LIVE!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == os.Getenv("BOT_ID") {
		return
	}

	if tripleCrown.TCActive && m.Author.ID == tripleCrown.TCUserId {
		tripleCrown.CollectMessages(s, m)
		return
	}

	command := strings.Fields(m.Content)[0]

	switch command {
	case "!gostatus":
		{
			status.BotStatus(s, m)
			break
		}
	case "!enditall":
		{
			i = 1
			s.ChannelMessageSend(m.ChannelID, "Goodbye!")
			break
		}
	case "!register":
		{
			facade.Register(s, m)
			break
		}
	case "!vouch":
		{
			vouch.Vouch(s, m)
			break
		}
	case "!unvouch":
		{
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			vouch.UnVouch(s, m)
			break
		}

	case "!repeat":
		{
			repeat.Repeat(s, m)
			break
		}
	case "!tcscore":
		{
			if tripleCrown.TCActive == false {
				embed := tripleCrown.TCrown(s, m)
				s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			}
		}
	}

}
