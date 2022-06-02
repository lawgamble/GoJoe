package vouch

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/facade/helpers"
	"os"
)

//TODO send message to whatever channel to let everyone know who's vouching for who.

func Vouch(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := vouchValidations(s, m)
	if err != nil {
		return
	}

	_ = giveUserVouchRole(s, m)

	successEmbed := embeds.CreateEmbed("VOUCHED!", "Vouch successful", "success")
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)

	s.ChannelMessageSend(os.Getenv("VOUCH_CHANNEL_ID"), fmt.Sprintf("<@%s>", m.Author.ID)+" has vouched for "+fmt.Sprintf("<@%s>", m.Mentions[0].ID))

	writeVouchedFile(m)
}

func vouchValidations(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !helpers.IsUserLeagueMember(m) {
		lmErr := fmt.Errorf("not League Member")
		return lmErr
	}

	errEmbed := discordgo.MessageEmbed{}

	err := validateMentions(m)

	if err != nil {
		switch error.Error(err) {
		case ">1":
			{
				errEmbed = embeds.CreateEmbed("NOPE", "You need to tag just ONE user", "danger")
				_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)
				return err
			}
		case "0":
			{
				errEmbed = embeds.CreateEmbed("NOPE", "You need to tag a user", "danger")
				_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)
				return err
			}
		case "dumbass":
			{
				_, _ = s.ChannelMessageSend(m.ChannelID, "https://tenor.com/view/gordon-ramsay-idiot-sandwich-angry-mad-what-are-you-gif-4169547")
				return err
			}
		}
	}
	return nil
}

func validateMentions(m *discordgo.MessageCreate) error {
	var err error
	if len(m.Mentions) > 1 {
		err = fmt.Errorf(">1")
		return err
	}
	if len(m.Mentions) == 0 || m.Mentions == nil {
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
