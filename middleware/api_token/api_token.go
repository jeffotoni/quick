package apitoken

import (
	"net/http"
)

const (
	errorAuthMessage = `{"error": "api-key is missing or it isn't correct", "code": 401}`
)

// Auth: This function gets the  API-Token provided by client from Header and check if it's valid
func Auth(apiTokenKey string, value string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		apiKeyFromHeader := req.Header.Get(apiTokenKey)
		if apiKeyFromHeader != value {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte(errorAuthMessage))
			return
		}
	}
}
