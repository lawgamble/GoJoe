package helpers

import "github.com/bwmarrin/discordgo"

func DoesMessageHaveMentions(m *discordgo.MessageCreate) bool {
	if len(m.Mentions) > 0 {
		return true
	}
	return false
}

//TODO get to this func when needed.
func DoesMessageHaveArguments(m *discordgo.MessageCreate) bool {
	// will need to check for m.Content but after the command itself
	return true
}

//TODO return message arguments in a slice if they are needed.
func MessageArguments(m *discordgo.MessageCreate) (args []string) {
	return args
}
