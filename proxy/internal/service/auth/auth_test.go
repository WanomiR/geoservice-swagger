package auth

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

var userAuth *UserAuth

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	userAuth = NewUserAuth("HT256", jwtSecret)
}

func TestNewUserAuth_RegisterUser(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"normal register request", User{"admin@example.com", "password"}, 201, "user registered"},
		{"wrong body", struct{ id int }{1}, 400, "EOF"},
		{"empty request", nil, 400, "EOF"},
		{"short password", User{"admin@example.com", "pa"}, 400, "password must be between 3 and 32 characters"},
		{"short email", User{"admi", "password"}, 400, "email must be between 5 and 32 characters"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/register", &body)
			wr := httptest.NewRecorder()

			userAuth.RegisterUser(wr, req)

			r := wr.Result()

			var resp JSONResponse
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

func TestNewUserAuth_Authenticate(t *testing.T) {
	testCases := []struct {
		name        string
		body        any
		wantStatus  int
		wantMessage string
	}{
		{"normal authenticate request", User{"admin@example.com", "password"}, 200, "user authenticated"},
		{"invalid credentials", User{"test@test.com", "password"}, 200, "invalid credentials"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tc.body)

			req := httptest.NewRequest("POST", "/api/login", &body)
			wr := httptest.NewRecorder()

			userAuth.Authenticate(wr, req)

			r := wr.Result()

			var resp JSONResponse
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
