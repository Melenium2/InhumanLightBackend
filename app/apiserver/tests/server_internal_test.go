package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/inhumanLightBackend/app/apiserver"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/inhumanLightBackend/app/utils/jwtHelper"
	"github.com/stretchr/testify/assert"
)

var (
	setAuthToken = func(r *http.Request) {
		jwt, _ := jwtHelper.Create(&models.User{
			ID: 1,
			Role: models.Roles[0],
		}, 1, "access")
		
		r.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", jwt))
	}
	httpParams = func (path string, method string, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
		rec := httptest.NewRecorder()
		bPayload := &bytes.Buffer{}
		if payload != nil {
			json.NewEncoder(bPayload).Encode(payload)
		}
		req, _ := http.NewRequest(method, path, bPayload)
		return rec, req
	}
)

func TestServer_HandleRegistration(t *testing.T) {
	server := apiserver.NewServer(teststore.New())
	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string {
				"email": "user123@gmail.com",
				"password": "123456",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string {
				"email": "user123@gmail.com",
				"password": "123",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bPayload := &bytes.Buffer{}
			json.NewEncoder(bPayload).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/signup", bPayload)
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleLogin(t *testing.T) {
	user := models.NewTestUser(t)
	email := "user123@gmail.com"
	password := "123456"
	user.Email = email
	user.Password = password

	store := teststore.New()
	assert.NoError(t, store.User().Create(user)) 

	server := apiserver.NewServer(store)
	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string {
				"email": email,
				"password": password,
			},
			expectedCode: http.StatusOK,
		}, 
		{
			name: "invalid body",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not find",
			payload: map[string]string {
				"email": email,
				"password": "123",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bPayload := &bytes.Buffer{}
			json.NewEncoder(bPayload).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/signin", bPayload)
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUserInfo(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)
	user.Role = models.Roles[1]
	server := apiserver.NewServer(store)

	testCases := []struct {
		name string
		tokenType string
		expectedCode int
	}{
		{
			name: "success",
			tokenType: "with token",
			expectedCode: http.StatusOK,
		},
		{
			name: "unauthorized",
			tokenType: "invalid token",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "success",
			tokenType: "without token",
			expectedCode: http.StatusUnauthorized,
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
			rec := httptest.NewRecorder()
			generatedAuth, _ := token(tc.tokenType)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/user?id=" + strconv.Itoa(user.ID), nil)
			req.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", generatedAuth))
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleRefreshAccessToken(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)
	server := apiserver.NewServer(store)

	testCases := []struct {
		name string
		tokenType string
		expectedCode int
	} {
		{
			name: "refresh token",
			tokenType: "valid",
			expectedCode: http.StatusOK,
		},
		{
			name: "access token",
			tokenType: "access",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			tokenType: "invalid",
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

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			generatedAuth, _ := token(tc.tokenType)
			req, _ := http.NewRequest(http.MethodGet, "/checkAccess", nil)
			req.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", generatedAuth))
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUpdateUser(t *testing.T) {
	user := models.NewTestUser(t)
	store := teststore.New()
	store.User().Create(user)
	server := apiserver.NewServer(store)

	testCases := []struct {
		name string
		payload interface{}
		expectedCode int
	} {
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
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			bPayload := &bytes.Buffer{}
			json.NewEncoder(bPayload).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/updateUser", bPayload)
			setAuthToken(req)
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)	
		})
	}
}

func TestServer_HandleTicketCreate(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	server := apiserver.NewServer(store)

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
			server.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestServer_HandleTicket(t *testing.T) {
	ticket := models.NewTestTicket(t)
	store := teststore.New()
	store.Tickets().Create(ticket)
	server := apiserver.NewServer(store)

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
			server.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}