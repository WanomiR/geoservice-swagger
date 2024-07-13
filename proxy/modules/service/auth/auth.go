package auth

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"proxy/modules/model"
	"proxy/modules/repository"
)

type Claims map[string]interface{}

type UserAuth struct {
	DB        repository.DatabaseRepo
	tokenAuth *jwtauth.JWTAuth
}

type UserAuthOption func(*UserAuth)

func WithDatabase(db repository.DatabaseRepo) UserAuthOption {
	return func(u *UserAuth) {
		u.DB = db
	}
}

func NewUserAuth(algorithm, secret string) *UserAuth {
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)

	userAuth := &UserAuth{
		tokenAuth: tokenAuth,
	}
	return userAuth
}

func (u *UserAuth) Register(user model.User) error {
	if len(user.Password) < 3 || len(user.Password) > 32 {
		return errors.New("password must be between 3 and 32 characters")
	}

	if len(user.Email) < 5 || len(user.Email) > 32 {
		return errors.New("email must be between 5 and 32 characters")
	}

	if _, err := u.DB.GetUserByEmail(user.Email); err == nil {
		return errors.New("user with this email already exists")
	}

	var newUser model.User
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	newUser.Password = string(encryptedPassword)
	newUser.Email = user.Email

	if err := u.DB.InsertUser(user); err != nil {
		return err
	}

	return nil
}

func (u *UserAuth) Authenticate(user model.User) (string, error) {
	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	err := bcrypt.CompareHashAndPassword([]byte(u.User.Password), []byte(user.Password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || u.User.Email != user.Email {
		w.WriteHeader(http.StatusOK)
		resp := JSONResponse{
			Code:    http.StatusOK,
			Message: "invalid credentials",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, tokenString, _ := u.tokenAuth.Encode(Claims{"email": user.Email})

	return tokenString, nil
}

func (u *UserAuth) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.VerifyRequest(u.tokenAuth, r, jwtauth.TokenFromCookie, jwtauth.TokenFromHeader)

		if err != nil || token == nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
