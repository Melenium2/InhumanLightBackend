package userroute

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/inhumanLightBackend/app/apiserver/apierrors"
	"github.com/inhumanLightBackend/app/apiserver/middleware"
	"github.com/inhumanLightBackend/app/apiserver/responses"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store"
)

type UserRoutes struct {
	store store.Store
}

func New(store store.Store) *UserRoutes {
	return &UserRoutes{
		store: store,
	}
}

func (ur *UserRoutes) SetUpRoutes(r *mux.Router) {
	r.HandleFunc("/user", ur.user()).Methods("GET")
	r.HandleFunc("/updateUser", ur.updateUser()).Methods("POST")
	r.HandleFunc("/notif/update", ur.updateNotif()).Methods("GET")
	r.HandleFunc("/notif/check", ur.checkNotif()).Methods("POST")
}

func (ur *UserRoutes) user() http.HandlerFunc {
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

		user, err := ur.store.User().FindById(userId)
		if err != nil {
			responses.SendError(w, r, http.StatusBadRequest, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, user)
	}
}

func (ur *UserRoutes) updateUser() http.HandlerFunc {
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
		authenticatedUser, err := ur.store.User().FindById(userId)
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

		if err := ur.store.User().Update(authenticatedUser); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "updated",
		})
	}
}

func (ur *UserRoutes) updateNotif() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxUser := middleware.UserContextMap(r.Context().Value(middleware.CtxUserKey))
		userId, err := strconv.Atoi(ctxUser["id"])
		if err != nil {
			responses.SendError(w, r, http.StatusUnauthorized, apierrors.ErrNotAuthenticated)
			return
		}

		notifs, err := ur.store.Notifications().FindById(uint(userId))
		if err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, notifs)
	}
}

func (ur *UserRoutes) checkNotif() http.HandlerFunc {
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

		if err := ur.store.Notifications().Check(req.Indexes, req.Id); err != nil {
			responses.SendError(w, r, http.StatusInternalServerError, err)
			return
		}

		responses.Respond(w, r, http.StatusOK, map[string]string{
			"message": "notifications updated",
		})
	}
}