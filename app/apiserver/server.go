package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func NewServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store: store,
	}

	s.configureRouter()
	
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/signup", handleRegistration(s)).Methods("POST")
	s.router.HandleFunc("/signin", handleLogin(s)).Methods("POST")
	s.router.HandleFunc("/checkAccess", handleRefreshAccessToken(s)).Methods("GET")
	prefix := s.router.PathPrefix("/api/v1").Subrouter()
	prefix.Use(authenticate)
	prefix.HandleFunc("/user", handleUserInfo(s)).Methods("GET")
	prefix.HandleFunc("/updateUser", handleUserUpdate(s)).Methods("POST")
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
