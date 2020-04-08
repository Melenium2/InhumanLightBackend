package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type server struct {
	router *mux.Router
	store  store.Store
}

func NewServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store: store,
	}

	s.configureRouter()
	
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(logging)
	s.router.HandleFunc("/signup", handleRegistration(s)).Methods("POST")
	s.router.HandleFunc("/signin", handleLogin(s)).Methods("POST")
	s.router.HandleFunc("/checkAccess", handleRefreshAccessToken(s)).Methods("GET")
	main := s.router.PathPrefix("/api/v1").Subrouter()
	main.Use(authenticate)
	main.HandleFunc("/user", handleUserInfo(s)).Methods("GET")
	main.HandleFunc("/updateUser", handleUserUpdate(s)).Methods("POST")
	support := main.PathPrefix("/support").Subrouter()
	support.HandleFunc("/ticket/create", handleCreateTicket(s)).Methods("POST")
	support.HandleFunc("/ticket", handleTicket(s)).Methods("GET")
	support.HandleFunc("/tickets", HandleTickets(s)).Methods("GET")
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func sendError(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string{"error": err.Error()})
}

func isAdmin(r *http.Request) bool {
	userCtx := userContextMap(r.Context().Value(ctxUserKey))
	return userCtx["access"] == models.Roles[1]
}
