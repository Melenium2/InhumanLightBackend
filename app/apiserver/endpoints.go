package apiserver

import (
	"encoding/json"
	"net/http"
)

func handleLogin() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		respond(w, r, 200, map[string]string{
			"response": "Hello",
		})
	}
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}