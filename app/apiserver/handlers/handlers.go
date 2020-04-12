package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	supportroutes "github.com/inhumanLightBackend/app/apiserver/handlers/supportroute"
	"github.com/inhumanLightBackend/app/apiserver/handlers/userroute"
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

	userroute.New(h.store).SetUpRoutes(main)
	supportroutes.New(h.store).SetUpRoutes(main)
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

