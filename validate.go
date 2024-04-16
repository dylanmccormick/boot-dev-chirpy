package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding json")
		w.WriteHeader(500)
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "too many characters")
	}
	cleanbody := removeProfanity(params.Body)

	chirp, err := cfg.db.CreateChirp(cleanbody)
	if err != nil {
		log.Fatal("unable to create chirp")
	}

	respondWithJSON(w, 200, chirp)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Fatal("Unable to get chirps", err)
	}
	type resp struct {
		Chirps []Chirp
	}

	response := resp{
		chirps,
	}

	respondWithJSON(w, 200, response)

}

func respondWithJSON(w http.ResponseWriter, code int, msg any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error encoding json")
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorBody struct {
		Error string `json:"error"`
	}

	resp := ErrorBody{msg}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error encoding json")
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func removeProfanity(s string) string {
	words := strings.Fields(s)
	for _, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle":
			s = strings.Replace(s, word, "****", 1)
		case "sharbert":
			s = strings.Replace(s, word, "****", 1)
		case "fornax":
			s = strings.Replace(s, word, "****", 1)
		}
	}

	return s
}
