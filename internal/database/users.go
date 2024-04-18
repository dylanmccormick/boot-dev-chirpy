package database

import (
	"errors"
	"fmt"
	"log"
)

type User struct {
	Email    string `json:"email"`
	Password []byte `'json:"password"`
	Id       int    `json:"id"`
}

func (db *DB) CreateUser(email string, password []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}
	id := len(data.Users) + 1
	user := User{email, password, id}
	log.Printf("Pre-assignment\n")

	data.Users[id] = user
	log.Printf("Writing user to db: %v\n", user)
	err = db.writeDB(data)

	if err != nil {
		return User{}, err
	}
	return user, nil

}

func (db *DB) GetUsers() ([]User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return []User{}, err
	}

	users := []User{}

	for _, v := range data.Users {
		users = append(users, v)
	}

	return users, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}

	for _, v := range users {
		if id == v.Id {
			return v, nil
		}
	}

	return User{}, errors.New(fmt.Sprintf("User with Id: %d not found", id))

}

func (db *DB) GetUser(email string) (User, error) {
	users, err := db.GetUsers()
	if err != nil {
		return User{}, err
	}

	for _, v := range users {
		if email == v.Email {
			return v, nil
		}
	}

	return User{}, errors.New(fmt.Sprintf("User with email: %s not found", email))

}

func (db *DB) UpdateUser(id int, email string, password []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.LoadDB()
	if err != nil {
		return User{}, err
	}

	user := data.Users[id]
	user.Password = password
	user.Email = email
	data.Users[id] = user

	err = db.writeDB(data)

	if err != nil {
		return User{}, err
	}
	return user, nil

}
