package main

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/zsolomon88/bootdev-chirpy/internal"
	"golang.org/x/crypto/bcrypt"
)

func createUserHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	userEmail := strings.TrimSpace(params.Email)
	userPassword := strings.TrimSpace(params.Password)

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, 500, "unable to create user")
		return
	}


	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Unable to connect to database")
		return
	}
	usr, err := dbHandle.CreateUser(userEmail, string(hashedPwd))
	if err != nil {
		respondWithError(w, 500, "Unable to write to database")
		return
	}
	type UserReply struct {
		Id int `json:"id"`
		Email string `json:"email"`
	}
	jsonReply := UserReply {
		Id: usr.Id,
		Email: usr.Email,
	}
	respondWithJSON(w, 201, jsonReply)
}

func authenticateHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	userEmail := strings.TrimSpace(params.Email)
	userPassword := strings.TrimSpace(params.Password)

	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Unable to connect to database")
		return
	}

	users, err := dbHandle.GetUsers()
	if err != nil {
		respondWithError(w, 500, "Unable to obtain data from db")
		return
	}

	for _, innerUser := range users {
		if innerUser.Email == userEmail {
			pwdMatch := bcrypt.CompareHashAndPassword([]byte(innerUser.Password), []byte(userPassword))
			if pwdMatch == nil {
				type UserResponse struct {
					Id int `json:"id"`
					Email string `json:"email"`
				}
				successResponse := UserResponse {
					Id: innerUser.Id,
					Email: innerUser.Email,
				}
				respondWithJSON(w, 200, successResponse)
				return
			}
		}
	}

	respondWithError(w, 401, "Unauthorized")
}
