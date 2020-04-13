package supportroutes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	"github.com/inhumanLightBackend/app/apiserver/middleware"
	"github.com/inhumanLightBackend/app/apiserver/responses"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type SupportRoutes struct {
	store store.Store
}

func New(store store.Store) *SupportRoutes {
	return &SupportRoutes{
		store: store,
	}
}

func (sr *SupportRoutes) SetUpRoutes(r *mux.Router) {
	support := r.PathPrefix("/support").Subrouter()

	support.HandleFunc("/ticket/create", sr.createTicket()).Methods("POST")
	support.HandleFunc("/ticket", sr.ticket()).Methods("GET")
	support.HandleFunc("/tickets", sr.tickets()).Methods("GET")
	support.HandleFunc("/message/add", sr.addMessage()).Methods("POST")
	support.HandleFunc("/messages", sr.messages()).Methods("GET")
	support.HandleFunc("/ticket/status", sr.changeMessageStatus()).Methods("GET")
}

func (sr *SupportRoutes) createTicket() http.HandlerFunc {
	type requset struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Section     string `json:"section"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requset{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		if req.Title == "" || req.Description == "" || req.Section == "" {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, _ := strconv.Atoi(ctxUser["id"])

		ticket := &models.Ticket{
			Title:       req.Title,
			Description: req.Description,
			Section:     req.Section,
			From:        uint(userId),
		}

		if err := sr.store.Tickets(r.Context()).Create(ticket); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "created",
		})
	}
}

func (sr *SupportRoutes) ticket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["id"]
		if !ok && len(id) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}
		ticketId, err := strconv.Atoi(id[0])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		ticket, err := sr.store.Tickets(r.Context()).Find(uint(ticketId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, ticket)
	}
}

func (sr *SupportRoutes) tickets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, err := strconv.Atoi(ctxUser["id"])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		tickets, err := sr.store.Tickets(r.Context()).FindAll(uint(userId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, tickets)
	}
}

func (sr *SupportRoutes) addMessage() http.HandlerFunc {
	type request struct {
		Message  string `json:"message"`
		TicketId uint	`json:"ticket_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		if req.Message == "" || req.TicketId == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, _ := strconv.Atoi(ctxUser["id"])
		if err := sr.store.Tickets(r.Context()).AddMessage(&models.TicketMessage{
			TicketId: req.TicketId,
			Message: req.Message,
			Who: uint(userId),
		}); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string {
			"message": "added",
		})
	}
}

func (sr *SupportRoutes) messages() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		id, ok := r.URL.Query()["id"]
		if !ok && len(id) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		ticketId, err := strconv.Atoi(id[0])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		messages, err := sr.store.Tickets(r.Context()).TakeMessages(uint(ticketId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, messages)
	}
}

func (sr *SupportRoutes) changeMessageStatus() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		id, ok := r.URL.Query()["id"]
		if !ok && len(id) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		status, ok := r.URL.Query()["st"]
		if !ok && len(status) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		ticketId, err := strconv.Atoi(id[0])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		if err := sr.store.Tickets(r.Context()).ChangeStatus(uint(ticketId), status[0]); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string {
			"message": "status chnaged to " + status[0],
		})
	}
}