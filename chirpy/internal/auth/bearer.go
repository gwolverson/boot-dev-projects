package auth

import (
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	apikey, err := StripAuthorizationHeader(headers, "Bearer")
	if err != nil {
		return "", err
	}
	return apikey, nil
}
