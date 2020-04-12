package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/inhumanLightBackend/app/apiserver/handlers"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/models/roles"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	setAuthToken = func(r *http.Request) {
		jwt, _ := jwtHelper.Create(&models.User{
			ID:   1,
			Role: roles.USER,
		}, 1, "access")

		r.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", jwt))
	}
	httpParams = func(path string, method string, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
		rec := httptest.NewRecorder()
		bPayload := &bytes.Buffer{}
		if payload != nil {
			json.NewEncoder(bPayload).Encode(payload)
		}
		req, _ := http.NewRequest(method, path, bPayload)
		return rec, req
	}
	newRequest = func(path string, method string, payload interface{}) *http.Request {
		bPayload := &bytes.Buffer{}
		if payload != nil {
			json.NewEncoder(bPayload).Encode(payload)
		}
		return httptest.NewRequest(method, path, bPayload)
	}
)

func TestServer_HandleSignUp(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name         string
		in           *http.Request
		out          *httptest.ResponseRecorder
		expectedCode int
	}{
		{
			name: "valid",
			in: newRequest("/signup", http.MethodPost, map[string]string{
				"email":    "user123@gmail.com",
				"password": "123456",
			}),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			in:           newRequest("/signup", http.MethodPost, "invalid"),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			in: newRequest("/signup", http.MethodPost, map[string]string{
				"email":    "user123@gmail.com",
				"password": "123",
			}),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h.SignUp().ServeHTTP(tc.out, tc.in)
			assert.Equal(t, tc.out.Code, tc.expectedCode)
		})
	}
}

func TestServer_HandleSignIn(t *testing.T) {
	email := "user123@gmail.com"
	password := "123456"

	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()
	user := models.NewTestUser(t)
	user.Email = email
	user.Password = password

	assert.NoError(t, store.User().Create(user))

	testCases := []struct {
		name         string
		in           *http.Request
		out          *httptest.ResponseRecorder
		expectedCode int
	}{
		{
			name: "valid",
			in: newRequest("/signin", http.MethodGet, map[string]string{
				"email":    email,
				"password": password,
			}),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid body",
			in:           newRequest("/signin", http.MethodGet, "invalid"),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not find",
			in: newRequest("/signin", http.MethodGet, map[string]string{
				"email":    email,
				"password": "123",
			}),
			out:          httptest.NewRecorder(),
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h.SignIn().ServeHTTP(tc.out, tc.in)
			assert.Equal(t, tc.expectedCode, tc.out.Code)
		})
	}
}

func TestServer_HandleRefreshAccessToken(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name         string
		in           *http.Request
		out          *httptest.ResponseRecorder
		tokenType    string
		expectedCode int
	}{
		{
			name:         "refresh token",
			in: newRequest("/checkAccess", http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType:    "valid",
			expectedCode: http.StatusOK,
		},
		{
			name:         "access token",
			in: newRequest("/checkAccess", http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType:    "access",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid token",
			in: newRequest("/checkAccess", http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType:    "invalid",
			expectedCode: http.StatusUnauthorized,
		},
	}

	token := func(tokenType string) (string, error) {
		switch tokenType {
		case "valid":
			return jwtHelper.Create(user, 30, "refresh")
		case "access":
			return jwtHelper.Create(user, 1, "access")
		case "invalid":
			jwt, err := jwtHelper.Create(user, 30, "refresh")

			return jwt + "123123", err
		}
		return "", nil
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generatedAuth, _ := token(tc.tokenType)
			tc.in.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", generatedAuth))
			h.ServeHTTP(tc.out, tc.in)
			assert.Equal(t, tc.expectedCode, tc.out.Code)
		})
	}
}

func TestServer_HandleUserInfo(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)
	user.Role = roles.ADMIN
	
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()
	
	testCases := []struct {
		name string
		in   *http.Request
		out  *httptest.ResponseRecorder
		tokenType string
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			in: newRequest("/api/v1/user?id=" + strconv.Itoa(user.ID), http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType: "with token",
			expectedCode: http.StatusOK,
			expectedBody: user.Login,
		},
		{
			name: "unauthorized",
			in: newRequest("/api/v1/user?id=" + strconv.Itoa(user.ID), http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType: "invalid token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "",
		},
		{
			name: "no success",
			in: newRequest("/api/v1/user?id=" + strconv.Itoa(user.ID), http.MethodGet, nil),
			out: httptest.NewRecorder(),
			tokenType: "without token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "",
		},
	}

		token := func (tokenType string) (string, error) {
			switch tokenType {
			case "with token":
				return jwtHelper.Create(user, 1, "access")
			case "invalid token":
				t, err := jwtHelper.Create(user, 1, "access")
				if err != nil {
					return "", err
				}

				return t + "123213", nil
			case "without token":
				return "", nil
			}

			return "", nil
		}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			generatedAuth, _ := token(tc.tokenType)
			tc.in.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", generatedAuth))
			h.ServeHTTP(tc.out, tc.in)
			assert.Equal(t, tc.expectedCode, tc.out.Code)
			response := &models.User{}
			assert.NoError(t, json.NewDecoder(tc.out.Body).Decode(response))
			assert.Equal(t, tc.expectedBody, response.Login)
		})
	}
}

func TestServer_HandleUpdateUser(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)

	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid email update",
			payload: map[string]string {
				"email": "123@user.com",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid contacts update",
			payload: map[string]string {
				"contacts": "CONTACTS",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid login update",
			payload: map[string]string {
				"login": "User_good_123",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Not admin trying to change ROLE",
			payload: map[string]string {
				"user_role": "ADMIN",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "Trying to change TOKEN",
			payload: map[string]string {
				"api_token": "1232131231",
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid email",
			payload: map[string]string {
				"email": "123@",
			},
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/updateUser", http.MethodPost, tc.payload)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleNotifUpdate(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name string
		authenticated bool
		expectedCode int
	}{
		{
			name: "valid",
			authenticated: true,
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid",
			authenticated: false,
			expectedCode: http.StatusUnauthorized,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/notif/update", http.MethodGet, nil)
			if tc.authenticated {
				setAuthToken(r)
			}
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleNotifCheck(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	for i := 0; i < 5; i++ {
		notif := models.NewTestNotification(t)
		store.Notifications().Create(notif)
	}

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{} {
				"id": 3,
				"indexes": []int{1, 2, 3, 4, 5},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "empty id",
			payload: map[string]interface{} {
				"indexes": []int{1, 2, 3, 4, 5},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty indexes",
			payload: map[string]interface{} {
				"id": 3,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid payload",
			payload: "text",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "negative user id",
			payload: map[string]interface{} {
				"id": -3,
				"indexes": []int{1, 2, 3, 4, 5},
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "negative indexes but valid",
			payload: map[string]interface{} {
				"id": 3,
				"indexes": []int{1, -2, 3, -4, 5},
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func (t *testing.T) {
			w, r := httpParams("/api/v1/notif/check", http.MethodPost, tc.payload)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}


func TestServer_HandleTicketCreate(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name string
		payload interface{}
		authenticated bool
		expectedCode int
	} {
		{
			name: "not authenticated",
			payload: map[string]string {
				"title": ticket.Title,
				"description": ticket.Description,
				"section": ticket.Section,
			},
			authenticated: false,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "valid",
			payload: map[string]string {
				"title": ticket.Title,
				"description": ticket.Description,
				"section": ticket.Section,
			},
			authenticated: true,
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid json",
			payload: "invalid",
			authenticated: true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty param",
			payload: map[string]string {
				"title": ticket.Title,
				"description": ticket.Description,
			},
			authenticated: true,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/support/ticket/create", http.MethodPost, tc.payload)
			if tc.authenticated {
				setAuthToken(r)
			}
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleTicket(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	store.Tickets().Create(ticket)
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name string
		path string
		expectedCode int
	} {
		{
			name: "valid",
			path: "1",
			expectedCode: http.StatusOK,
		},
		{
			name: "empty param",
			path: "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid",
			path: "invalid",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams(fmt.Sprintf("/api/v1/support/ticket?id=%s", tc.path), http.MethodGet, nil)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleTickets(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	ticketCount := 5
	for i := 0; i < ticketCount; i++ {
		ticket := models.NewTestTicket(t)
		store.Tickets().Create(ticket)
	}

	testCases := []struct {
		name string
		authenticated bool
		expectedCode int
	} {
		{
			name: "valid",
			authenticated: true,
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid",
			authenticated: false,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/support/tickets", http.MethodGet, nil)
			if tc.authenticated {
				setAuthToken(r)
			}
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleAddMessage(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	} {
		{
			name: "valid",
			payload: map[string]interface{} {
				"message": "Some message",
				"ticket_id": 4,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid. megative ticket_id",
			payload: map[string]interface{} {
				"message": "Some message",
				"ticket_id": -3,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not valid. no message",
			payload: map[string]interface{} {
				"ticket_id": 4,
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/support/message/add", http.MethodPost, tc.payload)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleTakeMessages(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()
	
	var ticketId uint = 3
	for i := 0; i < 5; i++ {
		message := models.NewTestTicketMessage(t)
		message.TicketId = ticketId
		store.Tickets().AddMessage(message)
	}

	testCases := []struct {
		name string
		path string
		expectedCode int
	}{
		{
			name: "valid",
			path: "?id=3",
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid. empty id",
			path: "?id=",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid. id doesn`t exist. empty response",
			path: "?id=555",
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid. id param doesn`t exist",
			path: "",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func (t *testing.T)  {
			w, r := httpParams("/api/v1/support/messages" + tc.path, http.MethodGet, nil)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleChangeStatus(t *testing.T) {
	store := teststore.New()
	h := handlers.New(store, logrus.New())
	h.SetupRoutes()

	ticket := models.NewTestTicket(t)
	store.Tickets().Create(ticket)

	testCases := []struct {
		name string
		path string
		expectedCode int
	} {
		{
			name: "valid",
			path: "?id=1&st=in process",
			expectedCode: http.StatusOK,
		},
		{
			name: "valid",
			path: "?id=1&st=closed",
			expectedCode: http.StatusOK,
		},
		{
			name: "empty id",
			path: "?id=&st=in process",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid id number",
			path: "?id=-30&st=in process",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid id. doesn`t exist",
			path: "?id=555&st=in process",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid status",
			path: "?id=1&st=in 123123process",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty id param",
			path: "?st=in process",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty status",
			path: "?id=1&st=",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "empty status param",
			path: "?id=1",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w, r := httpParams("/api/v1/support/ticket/status" + tc.path, http.MethodGet, nil)
			setAuthToken(r)
			h.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}



