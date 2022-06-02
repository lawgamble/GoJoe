package helpers

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

func IsUserLeagueMember(m *discordgo.MessageCreate) bool {
	userRoles := m.Member.Roles
	leagueMemberId := os.Getenv("LEAGUE_MEMBER_ID")

	for i := range userRoles {
		if userRoles[i] == leagueMemberId {
			return true
		}
	}
	return false
}
