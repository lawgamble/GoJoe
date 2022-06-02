package vouch

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/facade/helpers"
	"goJoe/internal/service"
	"io/ioutil"
	"os"
)

func UnVouch(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !helpers.DoesMessageHaveMentions(m) {
		errEmbed := embeds.CreateEmbed("NOPE", "You need to tag a user", "danger")

		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)

		return
	}
	if canUserUnvouch(m) {
		_ = s.GuildMemberRoleRemove(m.GuildID, m.Mentions[0].ID, os.Getenv("VOUCH_ROLE_ID"))

		successEmbed := embeds.CreateEmbed("UN-VOUCHED!", "Un-vouch successful", "success")

		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)

		s.ChannelMessageSend(os.Getenv("VOUCH_CHANNEL_ID"), fmt.Sprintf("<@%s>", m.Author.ID)+" has un-vouched for "+fmt.Sprintf("<@%s>", m.Mentions[0].ID))

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
