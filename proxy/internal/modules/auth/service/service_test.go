package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"proxy/internal/modules/auth/entities"
	"proxy/internal/modules/auth/service/mock_repository"
	"testing"
)

var mockUser = entities.User{
	Email:    "admin@example.com",
	Password: "password",
}

func TestUserAuth_Register(t *testing.T) {
	testCases := []struct {
		name    string
		user    entities.User
		wantErr error
	}{
		{"incorrect password", entities.User{"test@test.com", "ps"}, ErrorBadPassword},
		{"incorrect email", entities.User{"te", "password"}, ErrorBadEmail},
		{"existing user", mockUser, ErrorUserExists},
		{"successful case", entities.User{"test@test.com", "password"}, nil},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockDb := NewMockDb(controller)
	userAuth := NewUserAuth("HS256", "verysecret", mockDb)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := userAuth.Register(tc.user)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Register(%v) error = %v, wantErr %v", tc.user.Email, err, tc.wantErr)
			}
		})
	}
}

func TestUserAuth_Authenticate(t *testing.T) {
	testCases := []struct {
		name    string
		user    entities.User
		wantErr error
	}{
		{"non-existing user", entities.User{"test@test.com", "password"}, ErrorInvalidCredentials},
		{"invalid password", entities.User{mockUser.Email, "agoajgeoh"}, ErrorInvalidCredentials},
		{"successful case", mockUser, nil},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	mockDb := NewMockDb(controller)
	userAuth := NewUserAuth("HS256", "verysecret", mockDb)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := userAuth.Authenticate(tc.user)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Authenticate(%v) error = %v, wantErr %v", tc.user.Email, err, tc.wantErr)
			}
		})
	}
}

func NewMockDb(controller *gomock.Controller) *mock_repository.MockDatabaseRepo {
	mockDb := mock_repository.NewMockDatabaseRepo(controller)

	mockDb.EXPECT().GetUserByEmail(gomock.Any()).DoAndReturn(func(email string) (entities.User, error) {
		switch email {
		case mockUser.Email:
			encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(mockUser.Password), bcrypt.DefaultCost)
			return entities.User{mockUser.Email, string(encryptedPassword)}, nil
		default:
			return entities.User{}, ErrorUserNotFound
		}
	}).AnyTimes()

	mockDb.EXPECT().InsertUser(gomock.Any()).Return(nil).AnyTimes()

	return mockDb
}
