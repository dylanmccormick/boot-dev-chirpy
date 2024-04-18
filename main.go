package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dylanmccormick/boot-dev-chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	Id             int
	fileServerHits int
	db             *database.DB
	debug          bool
	JWT            string
}

func main() {

	dbg := flag.Bool("debug", false, "Enable Debug Mode")
	flag.Parse()

	startServer(dbg)

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

func startServer(debug *bool) {
	const filePathRoot = "."
	const port = "8080"

	apiConfig := &apiConfig{}
	apiConfig.debug = *debug
	var err error
	apiConfig.db, err = database.NewDB("database.json", apiConfig.debug)
	if err != nil {
		log.Fatal("Unable to create database")
	}
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	apiConfig.JWT = jwtSecret

	mux := http.NewServeMux()
	handler := http.FileServer(http.Dir("."))
	wrappedHandler := apiConfig.middlewareMetricsInc(handler)
	mux.Handle("/app/*", http.StripPrefix("/app", wrappedHandler))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("/api/reset", apiConfig.handleReset)
	mux.HandleFunc("POST /api/chirps", apiConfig.postChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.getChirp)
	mux.HandleFunc("POST /api/users", apiConfig.postUser)
	mux.HandleFunc("POST /api/login", apiConfig.postLogin)
	mux.HandleFunc("PUT /api/users", apiConfig.putUser)
	mux.HandleFunc("POST /api/refresh", apiConfig.refreshAccessToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.revokeRefreshToken)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}
