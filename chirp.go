package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	database "github.com/zsolomon88/bootdev-chirpy/internal"
)

func (cfg *apiConfig) deleteHandle(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	tokenParts := strings.Split(authToken, " ")

	if len(tokenParts) != 2 {
		respondWithError(w, 403, "Invalid token recieved")
		return
	}
	if tokenParts[0] != "Bearer" {
		respondWithError(w, 403, "Invalid token type")
		return
	}

	tokenStr := tokenParts[1]

	tokenClaims := &jwt.RegisteredClaims{}
	plainToken, err := jwt.ParseWithClaims(tokenStr, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, 403, fmt.Sprintf("Incorrect token: %s", err))
		return
	}
	usrId, err := plainToken.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 403, "Incorrect token")
		return
	}
	authorToDelete, err := strconv.Atoi(usrId)
	if err != nil {
		respondWithError(w, 403, "Invalid user id")
		return
	}

	chirpId := r.PathValue("chirpId")
	idToDelete, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, 500, "Error with chirp id")
		return
	}

	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "DB error")
		return
	}
	chirps, err := dbHandle.GetChirps()
	if err != nil {
		respondWithError(w, 500, "DB error")
		return
	}
	for _, chirp := range chirps {
		if chirp.Id == idToDelete && chirp.Author == authorToDelete {
			err = dbHandle.DeleteChirp(idToDelete)
			if err == nil {
				respondWithJSON(w, 204, "")
				return
			} else {
				respondWithError(w, 500, fmt.Sprintf("DB error: %v", err))
				return
			}
		}
	}

	respondWithError(w, 403, "Chirp not found")
}

func (cfg *apiConfig) createHandle(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Authorization")
	tokenParts := strings.Split(authToken, " ")

	if len(tokenParts) != 2 {
		respondWithError(w, 403, "Invalid token recieved")
		return
	}
	if tokenParts[0] != "Bearer" {
		respondWithError(w, 403, "Invalid token type")
		return
	}

	tokenStr := tokenParts[1]

	tokenClaims := &jwt.RegisteredClaims{}
	plainToken, err := jwt.ParseWithClaims(tokenStr, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Incorrect token: %s", err))
		return
	}
	usrId, err := plainToken.Claims.GetSubject()
	if err != nil {
		respondWithError(w, 401, "Incorrect token")
		return
	}
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
	authorId, _ := strconv.Atoi(usrId)
	chirp, err := dbHandle.CreateChirp(validBody.CleanBody, authorId)
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
	sortMethod := r.URL.Query().Get("sort")
	if sortMethod == "asc" || sortMethod == "" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id < chirps[j].Id
		})
	} else if sortMethod == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].Id > chirps[j].Id
		})
	}

	authorToGet := r.URL.Query().Get("author_id")
	if authorToGet != "" {
		authorFilter, _ := strconv.Atoi(authorToGet)
		filterAuthor := []database.Chirp{}
		for _, innerChirp := range chirps {
			if innerChirp.Author == authorFilter {
				filterAuthor = append(filterAuthor, innerChirp)
			}
		}
		chirps = filterAuthor
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
