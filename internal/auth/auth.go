package auth

import (
	"errors"
	"net/http"
	"strings"
)

const (
	AuthTypeBearer string = "Bearer"
	AuthTypeAPIKey string = "ApiKey"
)

func GetAuthToken(header http.Header, authType string) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization header")
	}

	prefixStr, tokenStr, found := strings.Cut(auth, " ")
	if !found || authType != prefixStr {
		return "", errors.New("malformed authorization header")
	}

	return tokenStr, nil
}
