package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"goJoe/internal/facade"
	"goJoe/internal/facade/repeat"
	"goJoe/internal/facade/status"
	"goJoe/internal/facade/tripleCrown"
	"goJoe/internal/facade/vouch"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
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
	case "tcteam":
		// read a file and check the value of flag
	}

}
