package handlers

import (
	"github.com/gorilla/mux"
)

type Route interface {
	SetUpRoutes(r *mux.Router)
}