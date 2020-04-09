package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/models/roles"
	"github.com/inhumanLightBackend/app/store"
)

// Server struct
type server struct {
	router *mux.Router
	store  store.Store
}

// Init new server
func NewServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store: store,
	}

	s.configureRouter()
	
	return s
}

// Handle requests
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Configure router and endpoints
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
	support.HandleFunc("/message/add", HandleAddMessage(s)).Methods("POST")
	support.HandleFunc("/messages", HandleMessages(s)).Methods("GET")
	support.HandleFunc("/ticket/status", HandleChangeStatus(s)).Methods("GET")
}

// Convert data in{} to the response body
func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Send error to the user
func sendError(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string{"error": err.Error()})
}

// Check if user in context have admin privilege
func isAdmin(r *http.Request) bool {
	userCtx := userContextMap(r.Context().Value(ctxUserKey))
	return userCtx["access"] == roles.ADMIN
}
