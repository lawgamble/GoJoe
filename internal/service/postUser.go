package service

import (
	"bytes"
	"net/http"
)

const postRequestURL = "http://45.77.195.189:8123/registerUser"

func Post(data []byte) *http.Response {
	res, err := http.Post(postRequestURL, "application/json", bytes.NewReader(data))
	if err != nil {
		// return the error here
	}
	if res.StatusCode != 200 {
		// send embed to tell user to try again later
	}
	return res
}
