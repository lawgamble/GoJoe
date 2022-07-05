package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const getUserURL = "http://45.77.195.189:8123/findPCLUserDiscordId/"

func GetUserData(userId string) (user UserResponse) {
	userURL := getUserURL + userId
	resp, err := http.Get(userURL)
	if err != nil {
		log.Println(err)
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		panic(err.Error())
	}
	_ = json.Unmarshal(body, &user)
	return user
}
