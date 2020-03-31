package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	
}

func newServer(store store.Store) *server {
	s := &server {
		router: mux.NewRouter(),
		logger: logrus.New(),
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/login", handleLogin())
	prefix := s.router.PathPrefix("/api/v1").Subrouter()
	prefix.Use(authenticate)
}
