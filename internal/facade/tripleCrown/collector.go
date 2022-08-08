package tripleCrown

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"log"
	"regexp"
	"strings"
)

func CollectMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.ToLower(m.Content) == "cancel" {
		chanResult = "cancel"
		result <- struct{}{}
		return
	}
	if len(collectionPool) == 0 {
		opponentId := extractId(m.Content)
		registered, opponent := validateRegistration(opponentId)
		if !registered {
			chanResult = "notRegistered"
			result <- struct{}{}
			return
		}
		Opponent = opponent.UserReg[0]
		uploadError := S.UploadUserData(Opponent)
		if uploadError != nil {
			errorMsg := embeds.CreateEmbed("Can't upload to DB", fmt.Sprintf("I had trouble uploading to the database: %v", uploadError), "gold")
			s.ChannelMessageSendEmbed(m.ChannelID, &errorMsg)
			return
		}
		collectionPool = append(collectionPool, Opponent.OculusName)
	} else {
		collectionPool = append(collectionPool, m.Content)
	}

	switch len(collectionPool) {
	case 1:
		fmt.Println(collectionPool[0])
		msg2 := embeds.CreateEmbed("Map", "What map did you play on?", "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg2)

	case 2:
		msg3 := embeds.CreateEmbed("Scores", "What was **your** score?", "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg3)

	case 3:
		msg4 := embeds.CreateEmbed("Scores", fmt.Sprintf("What was **%v's** score?", collectionPool[0]), "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg4)

	case 4:
		msg5 := embeds.CreateEmbed("Is this correct?", fmt.Sprintf("\nYou played: **%v**\nOn: **%v**\nYour score: **%v**\n%v's score: **%v**\n\nDoes this look accurate?\nType `yes` or `no`", collectionPool[0], collectionPool[1], collectionPool[2], collectionPool[0], collectionPool[3]), "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg5)

	case 5:
		if collectionPool[4] == "yes" {
			// do stuff with data for scoreboard embed
			resultEmbed := buildScoreResultEmbed(collectionPool)
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &resultEmbed)
			resetGlobalVars()
			chanResult = "complete"
			result <- struct{}{}
		} else {
			chanResult = "notCorrect"
			result <- struct{}{}
		}
	}
}

func extractId(content string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	id := reg.ReplaceAllString(content, "")

	return id
}
