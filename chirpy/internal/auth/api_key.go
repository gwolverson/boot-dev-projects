package auth

import (
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apikey, err := StripAuthorizationHeader(headers, "ApiKey")
	if err != nil {
		return "", err
	}
	return apikey, nil
}
