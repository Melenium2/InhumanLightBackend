package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
)

var (
	errIncorrectEmailOrPassword = errors.New("Incorrect email or password")
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
			sendError(w, r, http.StatusBadRequest, err)
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
			sendError(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := s.store.User().FindByEmail(req.Login)
		if err != nil || !user.ComparePassword(req.Password) {
			sendError(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		accToken, err := jwtHelper.CreateJwtToken(user, time.Now().Add(time.Hour*24).Unix())
		if err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return 
		}

		refrToken, err := jwtHelper.CreateJwtToken(user, time.Now().Add(time.Hour*24*30).Unix())
		if err != nil {
			sendError(w, r, http.StatusInternalServerError, err)
			return 
		}

		respond(w, r, http.StatusOK, map[string]string {
			"access_token": accToken,
			"refresh_token": refrToken,
		})
	}
}

// endpoint: api/v1//user?id=<?id>
func handleUserInfo(s *server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["id"]
		if !ok && len(id[0]) == 0 {
			sendError(w, r, http.StatusBadRequest, errors.New("Invalid id param"))
			return
		}

		userId, err := strconv.Atoi(id[0])
		if err != nil {
			sendError(w, r, http.StatusBadRequest, errors.New("Invalid id param"))
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


func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func sendError(w http.ResponseWriter, r *http.Request, code int, err error) {
	respond(w, r, code, map[string]string {"error": err.Error()})
}
