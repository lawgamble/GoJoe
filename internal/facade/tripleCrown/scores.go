package tripleCrown

import (
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"strconv"
)

func buildScoreResultEmbed(pool []string) discordgo.MessageEmbed {
	var winner string
	var winnerScore string
	var loserScore string
	player := TCUserName
	opponent := pool[0]
	playerScore, _ := strconv.Atoi(pool[2])
	opponentScore, _ := strconv.Atoi(pool[3])

	if playerScore > opponentScore {
		winner = player
		winnerScore = pool[2]
		loserScore = pool[3]
	} else {
		winner = opponent
		winnerScore = pool[3]
		loserScore = pool[2]
	}

	embed := embeds.ScoreResultsEmbed(pool, player, winner, winnerScore, loserScore, "gold")

	return embed
}
