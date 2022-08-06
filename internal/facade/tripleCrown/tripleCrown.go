package tripleCrown

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"runtime"
	"time"
)

var TCActive bool
var TCUserId string
var TCUserName string
var collectionPool []string
var chanResult string

var result chan struct{}

func TCrown(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	fmt.Println(runtime.NumGoroutine())
	if registered := validateRegistration(m.Author.ID); !registered {
		return embeds.NotRegisteredEmbed
	}

	setGlobalVars(m)
	sendUserFirstMessage(s, m)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// goroutine that cancels process when time is over.
	go func() {
		select {
		case <-ctx.Done():
			cancel()
			break
		}
	}()

	select {
	case <-result:
		if chanResult == "cancel" {
			resetGlobalVars()
			break
		}
		if chanResult == "complete" {
			resetGlobalVars()
			return embeds.SuccessEmbed
		}
		if chanResult == "notCorrect" {
			resetGlobalVars()
			return embeds.NotCorrect
		}

	case <-ctx.Done():
		resetGlobalVars()
		close(result)
		return embeds.TimeoutEmbed
	}
	cancelMsg := embeds.CreateEmbed("Cancelled", "Wait a minute, then try again", "danger")
	return cancelMsg
}

func sendUserFirstMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg1 := embeds.CreateEmbed("Opponent", "Who did you play?", "gold")
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg1)
}

func setGlobalVars(m *discordgo.MessageCreate) {
	TCActive = true
	TCUserId = m.Author.ID
	result = make(chan struct{}, 1)
}

func resetGlobalVars() {
	TCActive = false
	TCUserId = ""
	collectionPool = []string{}
}
