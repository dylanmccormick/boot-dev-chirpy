package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	Id             int
	fileServerHits int
	db             *DB
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	apiConfig := &apiConfig{}
	var err error
	apiConfig.db, err = NewDB("database.json")
	if err != nil {
		log.Fatal("Unable to create database")
	}

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	wrappedHandler := apiConfig.middlewareMetricsInc(handler)
	mux.Handle("/app/*", http.StripPrefix("/app", wrappedHandler))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("/api/reset", apiConfig.handleReset)
	mux.HandleFunc("POST /api/chirps", apiConfig.postChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirps)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html>\n\n<body>\n<h1>Welcome, Chirpy Admin</h1>\n"))
	w.Write([]byte(fmt.Sprintf("<p>Chirpy has been visited %d times!</p>\n</body>\n\n</html>", cfg.fileServerHits)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits += 1

		next.ServeHTTP(w, r)
	})
}
