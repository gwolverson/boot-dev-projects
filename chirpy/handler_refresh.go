package main

import (
	"net/http"
	"time"

	"github.com/gwolverson/go-courses/chirpy/internal/auth"
	"github.com/gwolverson/go-courses/chirpy/internal/database"
)

func (apiConfig *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type tokenResponse struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Only authenticated users can create chirps", nil)
		return
	}

	dbToken, err := apiConfig.queries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil || tokenIsExpired(dbToken) || !dbToken.RevokedAt.Time.IsZero() {
		respondWithError(w, http.StatusUnauthorized, "Refresh token is invalid or expired", nil)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, apiConfig.signingSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating user jwt", err)
		return
	}

	respondWithJSON(w, http.StatusOK, tokenResponse{
		Token: accessToken,
	})
}

func (apiConfig *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil || refreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Only authenticated users can create chirps", nil)
		return
	}

	err = apiConfig.queries.UpdateRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update refresh token", err)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func tokenIsExpired(dbToken database.RefreshToken) bool {
	tokenExpiry := dbToken.ExpiresAt
	sixtyDaysFromNow := time.Now().Add(60 * 24 * time.Hour)
	if tokenExpiry.After(sixtyDaysFromNow) {
		return true
	}
	return false
}
