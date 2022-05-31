package vouch

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"goJoe/internal/service"
	"io/ioutil"
)

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
