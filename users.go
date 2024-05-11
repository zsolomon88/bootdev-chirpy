package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	database "github.com/zsolomon88/bootdev-chirpy/internal"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUserHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
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
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	jsonReply := UserReply{
		Id:    usr.Id,
		Email: usr.Email,
	}
	respondWithJSON(w, 201, jsonReply)
}

func (cfg *apiConfig) updateUsrHandle(w http.ResponseWriter, r *http.Request) {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	intId, _ := strconv.Atoi(usrId)
	updatedUserInfo := database.User{
		Email:    userEmail,
		Password: string(hashedPwd),
		Id:       intId,
	}
	_, err = dbHandle.UpdateUser(intId, updatedUserInfo)
	if err != nil {
		respondWithError(w, 500, "Unable to write to database")
		return
	}

	type updateResponse struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}
	respondWithJSON(w, 200, updateResponse{Email: updatedUserInfo.Email, Id: updatedUserInfo.Id})
}

func (cfg *apiConfig) refreshHandle(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "Unable to open db")
		return
	}

	validToken, err := dbHandle.CheckRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Invalid refresh token: %v", err))
		return
	}


	currentTime := time.Now()
	expiration := time.Now().Add(time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(currentTime), ExpiresAt: jwt.NewNumericDate(expiration), Subject: fmt.Sprintf("%d", validToken.Id)})
	signedToken, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, 500, "Unable to create token")
		return
	}

	type RespToken struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, RespToken{Token: signedToken})
}

func (cfg *apiConfig) revokeTokenHandle(w http.ResponseWriter, r *http.Request) {
	dbHandle, err := database.NewDB("./database.json")
	if err != nil {
		respondWithError(w, 500, "unable to open db")
		return
	}

	refreshToken := r.Header.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	err = dbHandle.DeleteToken(refreshToken)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Invalid refresh token: %v", err))
		return
	}

	respondWithJSON(w, 204, "")
}

func (cfg *apiConfig) authenticateHandle(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Expiration int    `json:"expires_in_seconds"`
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

	expTime := time.Hour
	if params.Expiration > 0 && params.Expiration < 60*60*24 {
		expTime = time.Second * time.Duration(params.Expiration)
	}

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
				currentTime := time.Now()
				expiration := time.Now().Add(expTime)
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: jwt.NewNumericDate(currentTime), ExpiresAt: jwt.NewNumericDate(expiration), Subject: fmt.Sprintf("%d", innerUser.Id)})
				signedToken, err := token.SignedString([]byte(cfg.jwtSecret))
				if err != nil {
					respondWithError(w, 500, "Unable to create token")
					return
				}

				refreshExpiration := time.Now().Add(60 * 24 * time.Hour)
				refreshToken, err := dbHandle.CreateRefreshToken(refreshExpiration, innerUser.Id)
				if err != nil {
					respondWithError(w, 500, "Unable to create refresh token")
					return
				}
				type UserResponse struct {
					Id    int    `json:"id"`
					Email string `json:"email"`
					Token string `json:"token"`
					RefreshToken string `json:"refresh_token"`
				}
				successResponse := UserResponse{
					Id:    innerUser.Id,
					Email: innerUser.Email,
					Token: signedToken,
					RefreshToken: refreshToken.Token,
				}
				respondWithJSON(w, 200, successResponse)
				return
			}
		}
	}

	respondWithError(w, 401, "Unauthorized")
}
