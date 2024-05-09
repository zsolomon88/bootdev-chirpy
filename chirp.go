package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	database "github.com/zsolomon88/bootdev-chirpy/internal"
)

func createHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	chirpBody := strings.TrimSpace(params.Body)
	if len(chirpBody) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	for _, badWord := range getBadWords() {
		chirpBody = replaceWord(chirpBody, badWord, "****")
	}
	type validReturn struct {
		CleanBody string `json:"cleaned_body"`
	}
	validBody := validReturn{
		CleanBody: chirpBody,
	}

	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Unable to connect to database")
		return
	}
	chirp, err := dbHandle.CreateChirp(validBody.CleanBody)
	if err != nil {
		respondWithError(w, 500, "Unable to write to database")
		return
	}
	respondWithJSON(w, 201, chirp)
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Unable to connect to database")
		return
	}

	chirps, err := dbHandle.GetChirps()
	if err != nil {
		respondWithError(w, 500, "Unable to obtain data from db")
		return
	}

	chirpId := r.PathValue("chirpId")
	if chirpId == "" {
		respondWithJSON(w, 200, chirps)
	} else {
		id, err := strconv.Atoi(chirpId)
		if err != nil {
			respondWithError(w, 500, "Issue getting chirp id")
			return
		}
		for _, chirp := range chirps {
			if chirp.Id == id {
				respondWithJSON(w, 200, chirp)
				return
			}
		}
		respondWithError(w, 404, "Chrip not found")
	}
}
