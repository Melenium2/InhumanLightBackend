package apiserver

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
)

var (
	
)

// endpoint: /signup
func handleRegistration(s *server) http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			sendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		user := &models.User{
			Login:    req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(user); err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		respond(w, r, http.StatusCreated, map[string]string{"response": "user created"})
	}
}

// endpoint: /signin
func handleLogin(s *server) http.HandlerFunc {
	type request struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			sendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		user, err := s.store.User().FindByEmail(req.Login)
		if err != nil || !user.ComparePassword(req.Password) {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrIncorrectEmailOrPassword)
			return
		}

		accToken, err := jwtHelper.Create(user, 1, "access")
		if err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return
		}

		refrToken, err := jwtHelper.Create(user, 30, "refresh")
		if err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return
		}

		respond(w, r, http.StatusOK, map[string]string{
			"access_token":  accToken,
			"refresh_token": refrToken,
		})
	}
}

// endpoint: api/v1//user?id=<?id>
func handleUserInfo(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isAdmin(r) {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			sendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		userId, err := strconv.Atoi(id[0])
		if err != nil {
			sendError(w, r, http.StatusBadRequest, apierrors.ErrEmptyParam)
			return
		}

		user, err := s.store.User().FindById(userId)
		if err != nil {
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		respond(w, r, http.StatusOK, user)
	}
}

// endpoint: api/v1/updateUser
func handleUserUpdate(s *server) http.HandlerFunc {

	isZeroValue := func(x interface{}) bool {
		return x == reflect.Zero(reflect.TypeOf(x)).Interface()
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userModel := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(userModel); err != nil {
			sendError(w, r, http.StatusBadRequest, apierrors.ErrNotValidBody)
			return
		}

		if !isAdmin(r) && userModel.Role != "" {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		if userModel.Token != "" {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrPermissionDenied)
			return
		}

		userCtx := userContextMap(r.Context().Value(ctxUserKey))

		userId, _ := strconv.Atoi(userCtx["id"])
		authenticatedUser, err := s.store.User().FindById(userId)
		if err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
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
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().Update(authenticatedUser); err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return
		}

		respond(w, r, http.StatusOK, map[string]string {
			"message": "updated",
		})
	}
}

// endpoint: /checkAccess
func handleRefreshAccessToken(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			sendError(w, r, http.StatusUnauthorized, err)
			return
		}

		claims, err := jwtHelper.Validate(token)
		if err != nil || claims.Type == "access" {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		accessToken, err := jwtHelper.Create(&models.User{
			ID:   claims.UserId,
			Role: claims.Access,
		}, 1, "access")

		if err != nil {
			sendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		respond(w, r, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": token,
		})
	}
}
