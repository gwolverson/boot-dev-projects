package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/gwolverson/go-courses/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	signingSecret  string
	polkaKey       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Failed to obtain DB connection")
	}
	dbQueries := database.New(db)

	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		queries:        dbQueries,
		platform:       os.Getenv("PLATFORM"),
		signingSecret:  os.Getenv("SIGNING_SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	mux := http.NewServeMux()

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)

	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUsers)
	mux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUsers)

	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handlerDeleteChirpById)
	mux.HandleFunc("POST /api/refresh", apiConfig.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handlerRevokeToken)

	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlerPolkaWebhooks)

	mux.HandleFunc("POST /api/login", apiConfig.handlerLogin)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Printf("Listening for requests on: %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
