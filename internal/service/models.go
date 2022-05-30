package service

import "github.com/bwmarrin/discordgo"

type UserResponse struct {
	Found   bool      `json:"Found"`
	UserReg []UserReg `json:"UserReg"`
}
type UserReg struct {
	DiscordID   string `json:"DiscordId"`
	OculusName  string `json:"OculusName"`
	Pcl         string `json:"PCL"`
	PavlovName  string `json:"PavlovName"`
	Score       string `json:"Score"`
	ScoreStatus string `json:"ScoreStatus"`
}

// AllRoles struct contains the functions to get a Discord Role
type AllRoles struct {
	VouchRole *discordgo.Role
}

type VouchedUsers struct {
	UserWhoVouched string
	VouchedUser    string
}

type JSONData struct {
	Data []VouchedUsers
}
