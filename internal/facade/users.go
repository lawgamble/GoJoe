package facade

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/embeds"
	"goJoe/internal/service"
	"strings"
)

func getAuthorId(m *discordgo.MessageCreate) string {
	return m.Author.ID
}

func Register(s *discordgo.Session, m *discordgo.MessageCreate) {
	oculusName, err := validateOculusName(m)
	if err != nil {
		errEmbed := embeds.CreateEmbed("You need a name!", "Add a name after the `!register` command!", "danger")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &errEmbed)
		return
	}
	authorId := getAuthorId(m)
	userData := service.GetUserData(authorId)

	if userData.Found == true && userData.UserReg[0].Pcl == "Y" {

		if userData.UserReg[0].OculusName != oculusName {
			changeOculusNameRequest(userData, oculusName, m)
			successEmbed := embeds.CreateEmbed("Success!", fmt.Sprintf("You were already registered, so we updated your name from **%v** to **%v**", userData.UserReg[0].OculusName, oculusName), "success")
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)
		} else {
			successEmbed := embeds.CreateEmbed("No Change:", fmt.Sprintf("You were already registered as **%v**", oculusName), "success")
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)
		}
	}
	if userData.Found == false {
		newRegistration(m, oculusName)

		successEmbed := embeds.CreateEmbed("Success!:", fmt.Sprintf("You're now registered as **%v**", oculusName), "success")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)

	}

	if userData.Found == true && userData.UserReg[0].Pcl == "N" {
		newRegistration(m, oculusName)

		successEmbed := embeds.CreateEmbed("Success!:", fmt.Sprintf("You're now registered as **%v**", oculusName), "success")
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &successEmbed)
	}
	return
}

func newRegistration(m *discordgo.MessageCreate, oculusName string) {
	newUserRequest := createNewUserRequest(m, oculusName)
	data, _ := json.Marshal(newUserRequest)
	res := service.Post(data)
	_ = res.Body.Close()
}

func createNewUserRequest(m *discordgo.MessageCreate, oculusName string) service.RegisterUserRequest {
	user := service.RegisterUserRequest{
		RegOrg:      "PCL",
		DiscordID:   m.Author.ID,
		PavlovName:  "",
		DiscordName: "",
		OculusName:  oculusName,
	}
	return user
}

func changeOculusNameRequest(d service.UserResponse, oculusName string, m *discordgo.MessageCreate) {
	updatedRequestStruct := updateRequestStruct(d, oculusName, m)
	data, _ := json.Marshal(updatedRequestStruct)
	res := service.Post(data)
	_ = res.Body.Close()
}

func updateRequestStruct(d service.UserResponse, oculusName string, m *discordgo.MessageCreate) service.RegisterUserRequest {
	reqStruct := service.RegisterUserRequest{
		RegOrg:      "PCL",
		DiscordID:   d.UserReg[0].DiscordID,
		DiscordName: m.Author.Username,
		OculusName:  oculusName,
	}
	return reqStruct
}

func validateOculusName(m *discordgo.MessageCreate) (string, error) {
	if len(strings.Fields(m.Content)) > 1 {
		oculusNameSlice := strings.Fields(m.Content)[1:]
		oculusName := ""
		for _, word := range oculusNameSlice {
			oculusName = oculusName + " " + word
		}
		return strings.TrimSpace(oculusName), nil
	}
	return "", fmt.Errorf("BAD")
}
