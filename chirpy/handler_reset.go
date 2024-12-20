package main

import (
	"errors"
	"net/http"
	"sync/atomic"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = atomic.Int32{}

	if cfg.platform != "dev" {
		respondWithError(
			w,
			http.StatusForbidden,
			"This endpoint is only accessible in development",
			errors.New("This endpoint is only accessible in development"),
		)
		return
	}

	err := cfg.queries.Reset(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Failed to delete users",
			errors.New("Failed to delete users"),
		)
		return
	}
}
