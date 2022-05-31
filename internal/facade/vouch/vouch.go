package vouch

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"os"
)

//TODO send message to whatever channel to let everyone know who's vouching for who.

func Vouch(s *discordgo.Session, m *discordgo.MessageCreate) {
	errEmbed := discordgo.MessageEmbed{}
	userId := m.Author.ID
	err := validateMentions(m)
	if err != nil {
		switch error.Error(err) {
		case ">1":
			{
				errEmbed = embeds.CreateEmbed("NOPE", "You need to tag just ONE user", "danger")
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>", userId))
				_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)
			}
		case "0":
			{
				errEmbed = embeds.CreateEmbed("NOPE", "You need to tag a user", "danger")
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>", userId))
				_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)
			}
		case "dumbass":
			{
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>", userId))
				_, _ = s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/gordon-ramsay-idiot-sandwich-angry-mad-what-are-you-gif-4169547")
			}
		}
	} else {
		vouchErr := giveUserVouchRole(s, m)
		if vouchErr != nil {
			return
		}
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>", userId))
	successEmbed := embeds.CreateEmbed("VOUCHED!", "You've vouched for "+m.Mentions[0].Username+", \nso go forth and conquer!", "success")
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)
	writeVouchedFile(m)
}

func validateMentions(m *discordgo.MessageCreate) error {
	var err error
	if len(m.Mentions) > 1 {
		err = fmt.Errorf(">1")
		return err
	}
	if len(m.Mentions) == 0 {
		err = fmt.Errorf("0")
		return err
	}
	if m.Mentions[0].ID == m.Author.ID {
		err = fmt.Errorf("dumbass")
		return err
	}
	return nil
}

func giveUserVouchRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	err := s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, os.Getenv("VOUCH_ROLE_ID"))
	if err != nil {
		return err
	}
	return nil
}
