package vouch

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/service"
	"io/ioutil"
	"os"
)

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
