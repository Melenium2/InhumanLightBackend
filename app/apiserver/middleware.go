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
		header := r.Header.Get("Authentication")
		println(header)
		if header == "" {
			sendError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		splitedToken := strings.Split(header, " ")
		if len(splitedToken) != 2 {
			sendError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		token := splitedToken[1]
		if token == "" {
			sendError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		claims, err := jwtHelper.ValidateJwtToken(token)
		if err != nil {
			sendError(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		
		ctx := context.WithValue(r.Context(), ctxUserKey, claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}