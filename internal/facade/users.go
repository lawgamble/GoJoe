package facade

import (
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/service"
)

func getAuthorId(m *discordgo.MessageCreate) string {
	return m.Author.ID
}

func Register(m *discordgo.MessageCreate) string {
	authorId := getAuthorId(m)
	userData := service.GetUserData(authorId)
	if userData.Found {
		return "You are officially registered!"
	} else {
		return "You are not registered!"
	}
}
