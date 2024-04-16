package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
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

	respondWithJSON(w, 201, chirp)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Fatal("Unable to get chirps", err)
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, 200, chirps)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Fatal("Unable to get chirps", err)
	}
	if r.PathValue("id") == "" {
		respondWithError(w, 400, fmt.Sprintf("Expected an id value, recieved: %v", r.PathValue("id")))

	}

	requestValue, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal("Error converting path value")
	}
	found := false
	for _, chirp := range chirps {
		if chirp.Id == requestValue {
			respondWithJSON(w, 200, chirp)
			found = true
		}

	}
	if !found {
		respondWithError(w, 404, fmt.Sprintf("Chirp not found with id %d", requestValue))
	}

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
