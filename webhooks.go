package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	database "github.com/zsolomon88/bootdev-chirpy/internal"
)

func (cfg *apiConfig) redWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	apiKey = strings.TrimPrefix(apiKey, "ApiKey ")
	if apiKey != cfg.polkaKey {
		respondWithError(w, 401, "Invalid API Key")
		return
	}
	type dataStruct struct {
		UserId int `json:"user_id"`
	}
	type parameters struct {
		Event string     `json:"event"`
		Data  dataStruct `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if strings.TrimSpace(params.Event) == "user.upgraded" {
		dbHandle, err := database.NewDB("./database.json")
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("DB error: %v", err))
		}

		userMap, err := dbHandle.GetUsers()
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("DB error: %v", err))
		}

		userToUpdate := database.User{}
		for _, usr := range userMap {
			if usr.Id == params.Data.UserId {
				userToUpdate = usr
				break
			}
		}

		userToUpdate.RedStatus = true
		_, err = dbHandle.UpdateUser(params.Data.UserId, userToUpdate)
		if err != nil {
			respondWithError(w, 404, "User not found")
			return
		}
		respondWithJSON(w, 200, "")
		return
	}
	respondWithError(w, 204, "")
}
