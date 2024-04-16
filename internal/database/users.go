package database

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Email string `json:"email"`
	Id    int    `json:"id"`
}

func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}
	id := len(data.Users) + 1
	user := User{email, id}
	log.Printf("Pre-assignment\n")

	data.Users[id] = user
	log.Printf("Writing user to db: %v\n", user)
	err = db.writeDB(data)

	if err != nil {
		return User{}, err
	}
	return user, nil

}
