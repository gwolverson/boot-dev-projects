package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/gwolverson/go-courses/chirpy/internal/auth"
	"github.com/gwolverson/go-courses/chirpy/internal/database"
)

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func (apiConfig *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", errors.New("invalid chirp id"))
		return
	}

	chirp, err := apiConfig.queries.GetChirp(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found for supplied id", errors.New("no chirp found for supplied id"))
		return
	}

	respondWithJSON(w, http.StatusOK, mapDbChirpToChirpResponse(chirp))
}

func (apiConfig *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		userId, err := uuid.Parse(authorId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user id", errors.New("invalid user id"))
			return
		}
		chirps, err := apiConfig.queries.GetChirpsByUserId(r.Context(), userId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
		dbChirps = append(dbChirps, chirps...)
	} else {
		chirps, err := apiConfig.queries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
		dbChirps = append(dbChirps, chirps...)
	}

	foundChirps := []Chirp{}
	for _, chirp := range dbChirps {
		foundChirps = append(foundChirps, mapDbChirpToChirpResponse(chirp))
	}

	sort.Slice(foundChirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return foundChirps[i].CreatedAt.After(foundChirps[j].CreatedAt)
		}
		return foundChirps[i].CreatedAt.Before(foundChirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, foundChirps)
}

func (apiConfig *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	} else if chirpContainsProfanity(&params.Body) {
		replaceProfaneWordsInChirp(&params.Body)
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil || token == "" {
		respondWithError(w, http.StatusUnauthorized, "Only authenticated users can create chirps", nil)
		return
	}

	uuid, jwtErr := auth.ValidateJWT(token, apiConfig.signingSecret)
	if jwtErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Only authenticated users can create chirps", nil)
		return
	}

	createChirpsParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: uuid,
	}
	chirp, err := apiConfig.queries.CreateChirp(r.Context(), createChirpsParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to create chirp", nil)
		return
	}

	respondWithJSON(w, http.StatusCreated, mapDbChirpToChirpResponse(chirp))
}

func (apiConfig *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	chirpUuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil || accessToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, apiConfig.signingSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	chirp, err := apiConfig.queries.GetChirp(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp found for supplied id", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Only the original author of the chirp can delete it", err)
		return
	}

	err = apiConfig.queries.DeleteChirp(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unexpected error occurred", nil)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func chirpContainsProfanity(chirp *string) bool {
	for _, word := range profaneWords {
		if strings.Contains(strings.ToLower(*chirp), word) {
			return true
		}
	}
	return false
}

func replaceProfaneWordsInChirp(chirp *string) {
	chirpWords := strings.Split(*chirp, " ")
	for _, word := range chirpWords {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			*chirp = strings.Replace(*chirp, word, "****", -1)
		}
	}
}

func mapDbChirpToChirpResponse(dbChirp database.Chirp) Chirp {
	return Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}
}
