package tripleCrown

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/service"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var TCActive bool
var TCUserId string
var TCUserName string
var collectionPool []string

var wg sync.WaitGroup

var result chan string
var stopChan chan string
var outcome string

func TCrown(s *discordgo.Session, m *discordgo.MessageCreate) discordgo.MessageEmbed {
	result = make(chan string)
	fmt.Println(runtime.NumGoroutine())

	TCActive = true
	TCUserId = m.Author.ID
	isRegistered := validateRegistration(m.Author.ID)
	if isRegistered == false {
		return embeds.NotRegisteredEmbed
	}

	//send user first message
	msg1 := embeds.CreateEmbed("Opponent", "Who did you play?", "gold")
	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg1)

	//start timer
	startTimer := time.NewTimer(30 * time.Second)

	// wait on the user to finish
	go func() {
		for {
			select {
			default:
				finished := wait() // this will be the collector func that returns true when completed
				if finished == true {
					result <- "fire"
					outcome = "complete"
					stopChan <- "fire"
				}
			case <-stopChan:
				return
			}
		}
	}()

	select {
	case <-result:
		if outcome == "cancel" {
			startTimer.Stop()
			resetGlobalVars()
			close(result)
			break
		}
		if outcome == "complete" {
			fmt.Println("Complete was sent")
			startTimer.Stop()
			resetGlobalVars()
			close(result)
			return embeds.SuccessEmbed
		}
		if outcome == "notCorrect" {
			startTimer.Stop()
			resetGlobalVars()
			close(result)
			return embeds.NotCorrect
		}

	case <-startTimer.C:
		resetGlobalVars()
		close(result)
		result = nil
		return embeds.TimeoutEmbed

	}
	// only when user cancels do we end up here.
	cancelMsg := embeds.CreateEmbed("Cancelled", "Wait a minute, then try again", "danger")
	return cancelMsg
}

func validateRegistration(id string) bool {
	userData := service.GetUserData(id)
	if userData.Found == true {
		if userData.UserReg[0].Pcl == "Y" {
			if userData.UserReg[0].OculusName != "" {
				TCUserName = userData.UserReg[0].OculusName
				return true
			}
		}
	}
	return false
}

func wait() bool {
	wg.Add(5)
	wg.Wait()
	return true
}

func CollectMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.ToLower(m.Content) == "cancel" {
		result <- "fire"
		outcome = "cancel"
		fmt.Println(outcome)
		return
	}

	collectionPool = append(collectionPool, m.Content)

	switch len(collectionPool) {
	case 1:
		msg2 := embeds.CreateEmbed("Map", "What map did you play on?", "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg2)
		wg.Done()

	case 2:
		msg3 := embeds.CreateEmbed("Scores", "What was **your** score?", "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg3)
		wg.Done()

	case 3:
		msg4 := embeds.CreateEmbed("Scores", fmt.Sprintf("What was **%v's** score?", collectionPool[0]), "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg4)
		wg.Done()

	case 4:
		msg5 := embeds.CreateEmbed("Is this correct?", fmt.Sprintf("\nYou played: **%v**\nOn: **%v**\nYour score: **%v**\n%v's score: **%v**\n\nDoes this look accurate?\nType `yes` or `no`", collectionPool[0], collectionPool[1], collectionPool[2], collectionPool[0], collectionPool[3]), "gold")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &msg5)
		wg.Done()

	case 5:
		if collectionPool[4] == "yes" {
			wg.Done()
			// do stuff with data for scoreboard embed
			resultEmbed := buildScoreResultEmbed(collectionPool)
			s.ChannelMessageSendEmbed(m.ChannelID, &resultEmbed)
			resetGlobalVars()
		} else {
			if collectionPool[4] != "yes" {
				result <- "fire"
				outcome = "notCorrect"
			}
		}
	}
}

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

func resetGlobalVars() {
	TCActive = false
	TCUserId = ""
	collectionPool = []string{}
	wg = sync.WaitGroup{}
	outcome = ""
}

//msgs := [3] {}
//who did you play?
//Gamble
//what map
//What was your score
//10
//what was Gamble's score
//7
//Joe v Gamble
//10-7
//Joe Won.
// "Who did you play?", "What Map?", "What was **your** score?"
