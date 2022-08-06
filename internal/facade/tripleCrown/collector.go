package tripleCrown

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"strings"
)

func CollectMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.ToLower(m.Content) == "cancel" {
		chanResult = "cancel"
		result <- struct{}{}
		return
	}

	collectionPool = append(collectionPool, m.Content)

	switch len(collectionPool) {
	case 1:
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
