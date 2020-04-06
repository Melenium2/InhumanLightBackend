package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/inhumanLightBackend/app/apiserver"
	"github.com/inhumanLightBackend/app/models"
	"github.com/inhumanLightBackend/app/store/teststore"
	"github.com/stretchr/testify/assert"
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
	server := apiserver.NewServer(store)

	testCases := []struct {
		name string
		id int
		expectedCode int
	}{
		{
			name: "unauthorized",
			id: 1,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/user?id=" + strconv.Itoa(tc.id), nil)
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}