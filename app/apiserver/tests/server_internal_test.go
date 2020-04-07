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
			return jwtHelper.CreateJwtToken(user, 1, "access")
		case "invalid token":
			t, err := jwtHelper.CreateJwtToken(user, 1, "access")
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
			return jwtHelper.CreateJwtToken(user, 30, "refresh")
		case "access":
			return jwtHelper.CreateJwtToken(user, 1, "access")
		case "invalid":
			jwt, err := jwtHelper.CreateJwtToken(user, 30, "refresh")

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
			jwt, _ := jwtHelper.CreateJwtToken(user, 1, "access")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/updateUser", bPayload)
			req.Header.Set("Authentication", fmt.Sprintf("%s %s", "Bearer", jwt))
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)	
		})
	}
}