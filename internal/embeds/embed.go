package embeds

import "github.com/bwmarrin/discordgo"

func CreateEmbed(title string, description string, color string) discordgo.MessageEmbed {
	e := discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       getEmbedColor(color),
	}
	return e
}

func getEmbedColor(color string) int {
	if color == "success" {
		return 3066993
	}
	if color == "caution" {
		return 16776960
	}
	if color == "danger" {
		return 15158332
	}
	return 0
}
