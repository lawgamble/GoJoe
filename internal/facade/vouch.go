package facade

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/service"
	"io/ioutil"
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

func writeVouchedFile(m *discordgo.MessageCreate) {
	dataStruct := service.JSONData{}

	newVouch := service.VouchedUsers{
		UserWhoVouched: m.Author.ID,
		VouchedUser:    m.Mentions[0].ID,
	}

	jsonFile, _ := ioutil.ReadFile("vouched.json")
	_ = json.Unmarshal(jsonFile, &dataStruct)

	if validateDuplicates(dataStruct, newVouch) {
		return
	}

	dataStruct.Data = append(dataStruct.Data, newVouch)

	bytes, _ := json.Marshal(&dataStruct)
	_ = ioutil.WriteFile("vouched.json", bytes, 0644)
}

func validateDuplicates(d service.JSONData, newVouch service.VouchedUsers) bool {
	vUsers := d.Data
	for i := range vUsers {
		if vUsers[i] == newVouch {
			return true
		}
	}
	return false
}

func UnVouch(s *discordgo.Session, m *discordgo.MessageCreate) {
	if canUserUnvouch(m) {
		_ = s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, os.Getenv("VOUCH_ROLE_ID"))
		successEmbed := embeds.CreateEmbed("UN-VOUCHED!", "You've un-vouched for "+m.Mentions[0].Username+" - \nthis is probably for the best...", "success")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)
	} else {
		errorEmbed := embeds.CreateEmbed("NO-CAN-DO!", "You can't un-vouch for "+m.Mentions[0].Username+".\nYou don't have permission.", "danger")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errorEmbed)
	}

}

func canUserUnvouch(m *discordgo.MessageCreate) bool {
	for i := range m.Member.Roles {
		if m.Member.Roles[i] == os.Getenv("LEAGUE_MANAGER_ROLE_ID") || m.Member.Roles[i] == os.Getenv("MOD_ROLE_ID") {
			return true
		}
		dataStruct := service.JSONData{}

		jsonFile, _ := ioutil.ReadFile("vouched.json")
		_ = json.Unmarshal(jsonFile, &dataStruct)

		for i := range dataStruct.Data {
			if dataStruct.Data[i].UserWhoVouched == m.Author.ID && dataStruct.Data[i].VouchedUser == m.Mentions[0].ID {
				return true
			}
		}
	}
	return false
}

func giveUserVouchRole(s *discordgo.Session, m *discordgo.MessageCreate) error {
	err := s.GuildMemberRoleAdd(m.GuildID, m.Mentions[0].ID, os.Getenv("VOUCH_ROLE_ID"))
	if err != nil {
		return err
	}
	return nil
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
