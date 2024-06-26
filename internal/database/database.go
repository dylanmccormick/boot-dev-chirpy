package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp `json:"chirps"`
	Users         map[int]User
	RevokedTokens map[string]time.Time
}

// creates a new db file and initiates a new struct for the api to leverage
func NewDB(path string, debug bool) (*DB, error) {
	if debug {
		_, err := os.OpenFile("database.json", os.O_WRONLY, 0666)
		if err == nil {
			log.Printf("Database exists %s\n", "database.json")
			log.Printf("Deleting and remaking database for debug mode\n")
			os.Remove("database.json")
		}

	}

	db := &DB{
		"database.json",
		sync.RWMutex{},
	}

	err := db.EnsureDB()

	if err != nil {
		return &DB{}, err
	}

	return db, nil

}

// creates a new DB file if it doesn't exist
func (db *DB) EnsureDB() error {
	_, err := os.OpenFile(db.path, os.O_WRONLY, 0666)
	if err == nil {
		log.Printf("Database exists %s\n", db.path)
		return nil
	}
	if os.IsNotExist(err) {
		log.Printf("DB file does not exist %s\n", db.path)
		err = os.WriteFile(db.path, nil, 0666)
		if err != nil {
			return err
		}

	}

	return nil

}

// loads db into memory
func (db *DB) LoadDB() (DBStructure, error) {
	dat := DBStructure{
		make(map[int]Chirp),
		make(map[int]User),
		make(map[string]time.Time),
	}

	resp, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	if len(resp) == 0 {
		return dat, nil
	}
	err = json.Unmarshal(resp, &dat)
	if err != nil {
		return DBStructure{}, err
	}
	return dat, nil

}

// saves db to file
func (db *DB) writeDB(dbstructure DBStructure) error {

	err := db.EnsureDB()
	if err != nil {
		return err
	}

	dat, err := json.Marshal(dbstructure)

	os.WriteFile(db.path, dat, 0666)

	return nil

}
