package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func StripAuthorizationHeader(headers http.Header, headerKey string) (string, error) {
	authHeader := headers.Get("Authorization")

	splitToken := strings.Split(authHeader, fmt.Sprintf("%s ", headerKey))

	if len(splitToken) != 2 {
		return "", errors.New("no specified token found in authorization header")
	}
	return splitToken[1], nil
}
