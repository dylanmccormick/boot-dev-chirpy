package database

import "log"

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

// creates a chirp and saves it to file
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(data.Chirps) + 1
	log.Printf("pre create chirp\n")
	chirp := Chirp{body, id}
	log.Printf("Pre-assignment\n")

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
