package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-Server/model"
	"net/http"
)

var sessions = map[string]string{}

func sighUp() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user model.SignupRequestUser

			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				logger.Error("Decode json error")
			}
			fmt.Println(user)
			fmt.Println("Test start")

			fmt.Println("Test end")
		})
}
