package repeat

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func Repeat(s *discordgo.Session, m *discordgo.MessageCreate) {
	deleteSingleMessage(s, m)

	trimmedRes := trimResponse(m.Content)

	_, _ = s.ChannelMessageSend(m.ChannelID, trimmedRes)
}

func trimResponse(c string) string {

	if len(strings.Fields(c)) > 1 {
		resSlice := strings.Fields(c)[1:]
		finalRes := ""
		for _, word := range resSlice {
			finalRes = finalRes + " " + word
		}
		return strings.TrimSpace(finalRes)
	}
	return "You need to say something!"
}

func deleteSingleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
}
