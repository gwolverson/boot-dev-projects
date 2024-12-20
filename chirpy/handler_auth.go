package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gwolverson/go-courses/chirpy/internal/auth"
	"github.com/gwolverson/go-courses/chirpy/internal/database"
)

func (apiConfig *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Password == "" || params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "A valid email and password are required to login", err)
		return
	}

	user, err := apiConfig.queries.FindUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	hashCheckError := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if hashCheckError != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, apiConfig.signingSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating user jwt", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed creating refresh token", err)
		return
	}

	expiresAt := time.Now().Add(60 * 24 * time.Hour)
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	}
	apiConfig.queries.CreateRefreshToken(r.Context(), createRefreshTokenParams)

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
