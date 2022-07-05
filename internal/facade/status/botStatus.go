package status

import (
	"github.com/bwmarrin/discordgo"
)

func BotStatus(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "I'm alive")
}
