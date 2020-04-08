package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/inhumanLightBackend/app/models"
)

var (
	errEmptyTicketBody = errors.New("Empty ticket body")
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
			sendError(w, r, http.StatusBadRequest, errEmptyTicketBody)
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