package apiserver

import "net/http"

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		
		
		next.ServeHTTP(w, r)
	})
}