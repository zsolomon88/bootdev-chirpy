package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id     int    `json:"id"`
	Body   string `json:"body"`
	Author int    `json:"author_id"`
}

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	RedStatus bool   `json:"is_chirpy_red"`
}

type RefreshToken struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
	Id         int       `json:"id"`
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	newDb := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := newDb.ensureDB()
	if err != nil {
		return nil, err
	}

	return &newDb, nil
}

// Create a new user and add to the db
func (db *DB) CreateUser(email string, password string) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}
	newUser := User{
		Id:        len(users) + 1,
		Email:     email,
		Password:  password,
		RedStatus: false,
	}

	users = append(users, newUser)
	userMap := make(map[int]User)

	for i, innerUser := range users {
		if _, ok := userMap[i+1]; ok {
			return User{}, fmt.Errorf("user %s already exists", email)
		}
		userMap[i+1] = innerUser
	}
	structure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	structure.Users = userMap

	errWrite := db.writeDB(structure)
	if errWrite != nil {
		return User{}, errWrite
	}
	return newUser, nil
}

func (db *DB) CreateRefreshToken(expiration time.Time, usrId int) (RefreshToken, error) {

	randomData := make([]byte, 32)
	rand.Read(randomData)
	refreshToken := hex.EncodeToString(randomData)

	tokenStruct := RefreshToken{
		Token:      refreshToken,
		Expiration: expiration,
		Id:         usrId,
	}

	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	if dbStruct.RefreshTokens == nil {
		dbStruct.RefreshTokens = make(map[string]RefreshToken)
	}

	dbStruct.RefreshTokens[tokenStruct.Token] = tokenStruct
	err = db.writeDB(dbStruct)
	if err != nil {
		return RefreshToken{}, err
	}
	return tokenStruct, nil
}

func (db *DB) CheckRefreshToken(token string) (RefreshToken, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	if _, ok := dbStruct.RefreshTokens[token]; !ok {
		return RefreshToken{}, fmt.Errorf("refresh token not found")
	}

	if time.Now().After(dbStruct.RefreshTokens[token].Expiration) {
		return RefreshToken{}, fmt.Errorf("refresh token expired")
	}

	return dbStruct.RefreshTokens[token], nil
}

func (db *DB) DeleteToken(token string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := dbStruct.RefreshTokens[token]; !ok {
		return fmt.Errorf("refresh token not found")
	}

	delete(dbStruct.RefreshTokens, token)

	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, author int) (Chirp, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:     len(chirps) + 1,
		Body:   body,
		Author: author,
	}

	chirps = append(chirps, newChirp)
	chirpMap := make(map[int]Chirp)

	for i, innerChirp := range chirps {
		chirpMap[i+1] = innerChirp
	}
	structure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	structure.Chirps = chirpMap

	errWrite := db.writeDB(structure)
	if errWrite != nil {
		return Chirp{}, errWrite
	}
	return newChirp, nil
}

// Update Pwd returns all chirps in the database
func (db *DB) UpdateUser(usrId int, update User) (User, error) {

	userMap, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if _, ok := userMap.Users[usrId]; !ok {
		return User{}, fmt.Errorf("user %d not found", usrId)
	}

	updatedUser := userMap.Users[usrId]
	updatedUser.Password = update.Password
	updatedUser.Email = update.Email
	updatedUser.RedStatus = update.RedStatus
	userMap.Users[usrId] = updatedUser
	err = db.writeDB(userMap)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

// GetUsers returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	user := []User{}

	userMap, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Chirp map len: %v\n", len(chirpMap.Chirps))
	for _, usr := range userMap.Users {
		user = append(user, usr)
	}
	return user, nil
}

// GetChirps returns all chirps in the database
func (db *DB) DeleteChirp(id int) error {
	chirpMap, err := db.loadDB()
	if err != nil {
		return err
	}
	if _, ok := chirpMap.Chirps[id]; ok {
		delete(chirpMap.Chirps, id)
		err = db.writeDB(chirpMap)
		return err
	}
	return fmt.Errorf("chirp not found")
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}

	chirpMap, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("Chirp map len: %v\n", len(chirpMap.Chirps))
	for _, tweet := range chirpMap.Chirps {
		chirps = append(chirps, tweet)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		_, err := os.Create(db.path)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	f, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	structure := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	//fmt.Printf("Read: %s\n", f)
	decodeErr := json.Unmarshal(f, &structure)
	if decodeErr != nil {
		return DBStructure{}, err
	}
	//fmt.Printf("Read chirp: %s\n", structure.Chirps[1].Body)

	return structure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	//	fmt.Printf("Writing: %s\n", dat)

	writeErr := os.WriteFile(db.path, dat, 0644)
	if writeErr != nil {
		return writeErr
	}

	return nil
}
