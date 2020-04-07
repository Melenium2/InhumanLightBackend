package apiserver

import (
	"context"
	"net/http"
	"strings"

	"github.com/inhumanLightBackend/app/utils/jwtHelper"
)

type ctxKey int8

const (
	ctxUserKey ctxKey = iota
)

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		token, err := getAuthToken(r)
		if err != nil {
			sendError(w, r, http.StatusUnauthorized, err)
			return 
		}

		claims, err := jwtHelper.ValidateJwtToken(token)
		if err != nil || claims.Type == "refresh" {
			sendError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		} 
		
		ctx := context.WithValue(r.Context(), ctxUserKey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAuthToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authentication")
	if header == "" {
		return "", errNotAuthenticated
	}

	splitedToken := strings.Split(header, " ")
	if len(splitedToken) != 2 {
		return "", errNotAuthenticated
	}

	token := splitedToken[1]
	if token == "" {
		return "", errNotAuthenticated
	}

	return token, nil
}