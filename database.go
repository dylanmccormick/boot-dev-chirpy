package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

// creates a new db file and initiates a new struct for the api to leverage
func NewDB(path string) (*DB, error) {
	db := &DB{
		"database.json",
		sync.RWMutex{},
	}

	err := os.WriteFile(db.path, nil, 0666)
	if err != nil {
		return &DB{}, err
	}

	return db, nil

}

// creates a chirp and saves it to file
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(data.Chirps)
	chirp := Chirp{id, body}

	data.Chirps[id] = chirp
	log.Printf("Writing chirp to db: %v\n", chirp)
	err = db.writeDB(data)

	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil

}

// Returns all chirps in the DB
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return []Chirp{}, err
	}

	chirps := []Chirp{}
	for _, v := range data.Chirps {
		chirps = append(chirps, v)

	}
	return chirps, nil

}

// creates a new DB file if it doesn't exist
func (db *DB) EnsureDB() error {
	_, err := os.OpenFile(db.path, os.O_WRONLY, 0666)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		fmt.Println("DB file does not exist")
		err = os.WriteFile(db.path, nil, 0666)
		if err != nil {
			return err
		}

	}

	return nil

}

// loads db into memory
func (db *DB) LoadDB() (DBStructure, error) {
	dat := DBStructure{}

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
