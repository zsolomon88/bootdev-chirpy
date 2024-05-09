package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		Id:   len(chirps) + 1,
		Body: body,
	}

	chirps = append(chirps, newChirp)
	chirpMap := make(map[int]Chirp)

	for i, innerChirp := range chirps {
		chirpMap[i+1] = innerChirp
	}
	structure := DBStructure{
		Chirps: chirpMap,
	}
	errWrite := db.writeDB(structure)
	if errWrite != nil {
		return Chirp{}, errWrite
	}
	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}

	chirpMap, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Chirp map len: %v\n", len(chirpMap.Chirps))
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

	fmt.Printf("Read: %s\n", f)
	decodeErr := json.Unmarshal(f, &structure)
	if decodeErr != nil {
		return DBStructure{}, err
	}
	fmt.Printf("Read chirp: %s\n", structure.Chirps[1].Body)

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
	fmt.Printf("Writing: %s\n", dat)

	writeErr := os.WriteFile(db.path, dat, 0644)
	if writeErr != nil {
		return writeErr
	}

	return nil
}