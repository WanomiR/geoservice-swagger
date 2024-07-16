package controller

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"log"
	"net/http/httptest"
	"proxy/internal/modules/auth/controller/mock_service"
	"proxy/internal/modules/auth/entities"
	"proxy/internal/modules/auth/service"
	"proxy/internal/utils/readresponder"
	"reflect"
	"testing"
)

var mockUser = entities.User{
	Email:    "admin@example.com",
	Password: "password",
}

func TestAuth_Register(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"successful registry", entities.User{"some@user.com", "password"}, 201, "user registered"},
		{"user exists", mockUser, 400, service.ErrorUserExists.Error()},
		{"wrong body", struct{ id int }{1}, 400, service.ErrorEOF.Error()},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)
	auth := NewAuth(mockService, readresponder.NewReadRespond())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/register", &body)
			wr := httptest.NewRecorder()

			auth.Register(wr, req)

			r := wr.Result()

			var resp readresponder.JSONResponse
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				log.Println(err)
			}

			defer r.Body.Close()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if tc.wantMessage != resp.Message {
				t.Errorf("got message %s, want %s", resp.Message, tc.wantMessage)
			}

		})
	}
}

func TestAuth_Authenticate(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"successful authentication", mockUser, 200, "user authenticated"},
		{"invalid credentials", entities.User{"test@test.com", "password"}, 400, service.ErrorInvalidCredentials.Error()},
		{"wrong body", struct{ id int }{1}, 400, service.ErrorEOF.Error()},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := NewMockService(controller)
	auth := NewAuth(mockService, readresponder.NewReadRespond())

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/register", &body)
			wr := httptest.NewRecorder()

			auth.Authenticate(wr, req)

			r := wr.Result()

			var resp readresponder.JSONResponse
			err := json.NewDecoder(r.Body).Decode(&resp)
			if err != nil {
				log.Println(err)
			}

			defer r.Body.Close()

			if r.StatusCode != tc.wantStatus {
				t.Errorf("got status code %d, want %d", r.StatusCode, tc.wantStatus)
			}

			if tc.wantMessage != resp.Message {
				t.Errorf("got message %s, want %s", resp.Message, tc.wantMessage)
			}

		})
	}
}

func NewMockService(controller *gomock.Controller) *mock_service.MockAuthenticator {
	mockService := mock_service.NewMockAuthenticator(controller)

	mockService.EXPECT().Register(gomock.Any()).DoAndReturn(func(user entities.User) error {
		switch {
		case reflect.DeepEqual(user, mockUser):
			return service.ErrorUserExists
		default:
			return nil
		}
	}).AnyTimes()

	mockService.EXPECT().Authenticate(gomock.Any()).DoAndReturn(func(user entities.User) (string, error) {
		switch {
		case reflect.DeepEqual(user, mockUser):
			return "token", nil
		default:
			return "", service.ErrorInvalidCredentials
		}
	}).AnyTimes()

	return mockService
}
