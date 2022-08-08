package tripleCrown

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/service"
	"time"
)

var TCActive bool
var TCUserId string
var TCUserName string
var collectionPool []string
var chanResult string

var result chan struct{}

var Challenger service.UserReg
var Opponent service.UserReg

var S Service

func TCrown(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	db, err := initDB()
	if err != nil {
		errorMsg := embeds.CreateEmbed("Can't Connect to DB", "I had trouble connecting to the database. Try again, bro.", "gold")
		return errorMsg
	}
	S.db = db
	defer db.Close()

	registered, userData := validateRegistration(m.Author.ID)
	if !registered {
		return embeds.NotRegisteredEmbed
	}
	uploadError := S.UploadUserData(userData.UserReg[0])
	if uploadError != nil {
		errorMsg := embeds.CreateEmbed("Can't upload to DB", fmt.Sprintf("I had trouble uploading to the database: %v", uploadError), "gold")
		return errorMsg
	}

	setGlobalVars(m, userData)
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
		if chanResult == "notRegistered" {
			resetGlobalVars()
			return embeds.OpponentNotRegisteredEmbed
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

func initDB() (*sql.DB, error) {
	db, err := startDB()
	if err != nil {
		return db, err
	}
	return db, nil
}

func sendUserFirstMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg1 := embeds.CreateEmbed("Opponent", "Who did you play?", "gold")
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg1)
}

func setGlobalVars(m *discordgo.MessageCreate, c service.UserResponse) {
	Challenger = c.UserReg[0]
	TCUserName = c.UserReg[0].OculusName
	TCActive = true
	TCUserId = m.Author.ID
	result = make(chan struct{}, 1)
}

func resetGlobalVars() {
	Challenger = service.UserReg{}
	TCUserName = ""
	TCUserId = ""
	TCActive = false
	collectionPool = []string{}
}
