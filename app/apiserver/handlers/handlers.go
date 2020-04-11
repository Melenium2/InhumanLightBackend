package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	"github.com/inhumanLightBackend/app/apiserver/middleware"
	"github.com/inhumanLightBackend/app/apiserver/responses"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	store  store.Store
	logger *logrus.Logger
	router *mux.Router
}

func New(store store.Store, logger *logrus.Logger) *Handlers {
	return &Handlers{
		store:  store,
		logger: logger,
		router: mux.NewRouter(),
	}
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handlers) SetupRoutes() {
	middleware := middleware.New(h.logger)
	h.router.Use(middleware.Logging)
	h.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	h.router.HandleFunc("/signup", h.SignUp()).Methods("POST")
	h.router.HandleFunc("/signin", h.SignIn()).Methods("POST")
	h.router.HandleFunc("/checkAccess", h.CheckAccessToken()).Methods("GET")

	main := h.router.PathPrefix("/api/v1").Subrouter()
	main.Use(middleware.Authenticate)

	main.HandleFunc("/user", h.User()).Methods("GET")
	main.HandleFunc("/updateUser", h.UpdateUser()).Methods("POST")
	main.HandleFunc("/notif/update", h.UpdateNotif()).Methods("GET")
	main.HandleFunc("/notif/check", h.CheckNotif()).Methods("POST")

	support := main.PathPrefix("/support").Subrouter()

	support.HandleFunc("/ticket/create", h.CreateTicket()).Methods("POST")
	support.HandleFunc("/ticket", h.Ticket()).Methods("GET")
	support.HandleFunc("/tickets", h.Tickets()).Methods("GET")
	support.HandleFunc("/message/add", h.AddMessage()).Methods("POST")
	support.HandleFunc("/messages", h.Messages()).Methods("GET")
	support.HandleFunc("/ticket/status", h.ChangeMessageStatus()).Methods("GET")
}

func (h *Handlers) SignUp() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		user := &models.User{
			Login:    req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := h.store.User().Create(user); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusCreated, map[string]string{"response": "user created"})
	}
}

func (h *Handlers) SignIn() http.HandlerFunc {
	type request struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		user, err := h.store.User().FindByEmail(req.Login)
		if err != nil || !user.ComparePassword(req.Password) {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrIncorrectEmailOrPassword)
			return
		}

		accToken, err := jwtHelper.Create(user, 1, "access")
		if err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		refrToken, err := jwtHelper.Create(user, 30, "refresh")
		if err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"access_token":  accToken,
			"refresh_token": refrToken,
		})
	}
}

func (h *Handlers) CheckAccessToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := middleware.GetToken(r)
		if err != nil {
			responses.SendError(w, r, http.StatusUnauthorized, err)
			return
		}

		claims, err := jwtHelper.Validate(token)
		if err != nil || claims.Type == "access" {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		accessToken, err := jwtHelper.Create(&models.User{
			ID:   claims.UserId,
			Role: claims.Access,
		}, 1, "access")

		if err != nil {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": token,
		})
	}
}

// MAIN ROUTE ////////////////////////////////////////////////////////////////////////////////////

func (h *Handlers) User() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !middleware.IsAdmin(r) {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		userId, err := strconv.Atoi(id[0])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		user, err := h.store.User().FindById(userId)
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, user)
	}
}

func (h *Handlers) UpdateUser() http.HandlerFunc {
	isZeroValue := func(x interface{}) bool {
		return x == reflect.Zero(reflect.TypeOf(x)).Interface()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userModel := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(userModel); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		if !middleware.IsAdmin(r) && userModel.Role != "" {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		if userModel.Token != "" {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		userCtx := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, _ := strconv.Atoi(userCtx["id"])
		authenticatedUser, err := h.store.User().FindById(userId)
		if err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		rNewModel := reflect.ValueOf(userModel)
		if rNewModel.Kind() == reflect.Ptr {
			rNewModel = rNewModel.Elem()
			for i := 0; i < rNewModel.NumField(); i++ {
				field := rNewModel.Field(i)
				if !isZeroValue(field.Interface()) {
					reflect.ValueOf(authenticatedUser).Elem().Field(i).Set(field)
				}
			}
		}

		if err := authenticatedUser.Validate(); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.store.User().Update(authenticatedUser); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "updated",
		})
	}
}

func (h *Handlers) UpdateNotif() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, err := strconv.Atoi(ctxUser["id"])
		if err != nil {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		notifs, err := h.store.Notifications().FindById(uint(userId))
		if err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, notifs)
	}
}

func (h *Handlers) CheckNotif() http.HandlerFunc {
	type request struct {
		Id      uint  `json:"id"`
		Indexes []int `json:"indexes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}
		if req.Id == 0 || len(req.Indexes) == 0 {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		if err := h.store.Notifications().Check(req.Indexes, req.Id); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "notifications updated",
		})
	}
}

// SUPPORT ROUTES //////////////////////////////////////////////////////////////////////////////////////

func (h *Handlers) CreateTicket() http.HandlerFunc {
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

		if err := h.store.Tickets().Create(ticket); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "created",
		})
	}
}

func (h *Handlers) Ticket() http.HandlerFunc {
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

		ticket, err := h.store.Tickets().Find(uint(ticketId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, ticket)
	}
}

func (h *Handlers) Tickets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, err := strconv.Atoi(ctxUser["id"])
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		tickets, err := h.store.Tickets().FindAll(uint(userId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, tickets)
	}
}

func (h *Handlers) AddMessage() http.HandlerFunc {
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
		if err := h.store.Tickets().AddMessage(&models.TicketMessage{
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

func (h *Handlers) Messages() http.HandlerFunc {
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

		messages, err := h.store.Tickets().TakeMessages(uint(ticketId))
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, messages)
	}
}

func (h *Handlers) ChangeMessageStatus() http.HandlerFunc {
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

		if err := h.store.Tickets().ChangeStatus(uint(ticketId), status[0]); err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string {
			"message": "status chnaged to " + status[0],
		})
	}
}
