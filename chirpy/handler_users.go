package main

import (
	"encoding/json"
	"net/http"

	"github.com/gwolverson/go-courses/chirpy/internal/auth"
	"github.com/gwolverson/go-courses/chirpy/internal/database"
)

type parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (apiConfig *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unexpected error", err)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	user, err := apiConfig.queries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (apiConfig *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil || accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", nil)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, apiConfig.signingSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unexpected error", err)
		return
	}

	updateUserParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	}

	user, err := apiConfig.queries.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
