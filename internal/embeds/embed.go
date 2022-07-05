package embeds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var NotRegisteredEmbed = CreateEmbed("You need to register!", "Use the `!register <NAME>` command to register with PCL.", "danger")
var TimeoutEmbed = CreateEmbed("Process Timeout", "You took too long! Try again.", "danger")
var SuccessEmbed = CreateEmbed("Process Complete", "Success!", "success")
var NotCorrect = CreateEmbed("Try Again", "You can go ahead and try again.\nBe honest and accurate!", "danger")

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
	if color == "gold" {
		return 15844367
	}
	return 0
}

func ScoreResultsEmbed(a []string, p string, w string, ws string, ls string, c string) discordgo.MessageEmbed {
	title := fmt.Sprintf("%v vs %v", p, a[0])
	description := fmt.Sprintf("%v Wins\n**%v-%v**", w, ws, ls)
	e := discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       getEmbedColor(c),
		Footer:      nil,
		Image:       nil,
		Thumbnail:   nil,
		Video:       nil,
		Provider:    nil,
		Author:      nil,
		Fields:      nil,
	}
	return e
}
