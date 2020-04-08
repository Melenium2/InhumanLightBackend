package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/inhumanLightBackend/app/models"
)

var (
	errEmptyParam = errors.New("Empty param")
)

// endpoint: api/v1/support/ticket/create
func handleCreateTicket(s *server) http.HandlerFunc  {
	type requset struct {
		Title string `json:"title"`
		Description string `json:"description"`
		Section string `json:"section"`
	}

	return func (w http.ResponseWriter, r *http.Request)  {
		req := &requset{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Title == "" || req.Description == "" || req.Section == "" {
			sendError(w, r, http.StatusBadRequest, errEmptyParam)
			return
		}

		ctxUser := userContextMap(r.Context().Value(ctxUserKey))
		userId, _ := strconv.Atoi(ctxUser["id"])

		ticket := &models.Ticket{
			Title: req.Title,
			Description: req.Description,
			Section: req.Section,
			From: uint(userId),
		}

		if err := s.store.Tickets().Create(ticket); err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return
		}

		respond(w, r, http.StatusOK, map[string]string {
			"message": "created",
		})
	}
}

// endpoint: api/v1/support/ticket?id=<?id>
// Return ticket by id
func handleTicket(s *server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		id, ok := r.URL.Query()["id"]
		if !ok && len(id) == 0 {
			sendError(w, r, http.StatusBadRequest, errEmptyParam)
			return
		}
		ticketId, err := strconv.Atoi(id[0])
		if err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}
		
		ticket, err := s.store.Tickets().Find(uint(ticketId))
		if err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		respond(w, r, http.StatusOK, ticket)
	}
}

// or endpoint: api/v1/support/tickets
// Return all сtxUser tickets 
func HandleTickets(s *server) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request)  {
		ctxUser := userContextMap(r.Context().Value(ctxUserKey))
		userId, err := strconv.Atoi(ctxUser["id"])
		if err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		tickets, err := s.store.Tickets().FindAll(uint(userId))
		if err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		respond(w, r, http.StatusOK, tickets)
	}
}