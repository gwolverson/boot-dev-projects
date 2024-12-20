package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gwolverson/go-courses/chirpy/internal/auth"
	"net/http"
)

const upgradedUserEvent = "user.upgraded"

type Webhook struct {
	Event string      `json:"event"`
	Data  WebhookData `json:"data"`
}

type WebhookData struct {
	UserId uuid.UUID `json:"user_id"`
}

func (apiConfig *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKeyIsInvalid(apiKey, apiConfig.polkaKey) {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	webhook := Webhook{}
	err = decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if webhook.Event != upgradedUserEvent {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	err = apiConfig.queries.UpdateUserToChirpyRed(r.Context(), webhook.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No user exists for supplied ID", err)
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}

func apiKeyIsInvalid(key string, polkaKey string) bool {
	return key == "" || key != polkaKey
}
