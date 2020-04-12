package balanceroute

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/apiserver/responses"
	"github.com/inhumanLightBackend/app/store"
)

type BalanceRoute struct {
	store store.Store
}

func New(store store.Store) *BalanceRoute {
	return &BalanceRoute{
		store: store,
	}
}

func (br *BalanceRoute) SetUp(r *mux.Router) {
	balance := r.PathPrefix("/balance").Subrouter()
	balance.HandleFunc("/check", br.check()).Methods("GET")
}

func (br *BalanceRoute) check() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		responses.Respond(w, r, http.StatusOK, "kalss")
	}
}