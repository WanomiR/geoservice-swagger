package auth

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"proxy/internal/entities"
	"proxy/internal/repository"
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

func NewUserAuth(algorithm, secret string, options ...UserAuthOption) *UserAuth {
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)

	userAuth := &UserAuth{
		tokenAuth: tokenAuth,
	}

	for _, option := range options {
		option(userAuth)
	}

	return userAuth
}

func (a *UserAuth) Register(user entities.User) error {
	if len(user.Password) < 3 || len(user.Password) > 32 {
		return errors.New("password must be between 3 and 32 characters")
	}

	if len(user.Email) < 5 || len(user.Email) > 32 {
		return errors.New("email must be between 5 and 32 characters")
	}

	if _, err := a.DB.GetUserByEmail(user.Email); err == nil {
		return errors.New("user with this email already exists")
	}

	var newUser entities.User
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	newUser.Password = string(encryptedPassword)
	newUser.Email = user.Email

	if err := a.DB.InsertUser(newUser); err != nil {
		return err
	}

	return nil
}

func (a *UserAuth) Authenticate(userQuery entities.User) (string, error) {
	user, err := a.DB.GetUserByEmail(userQuery.Email)
	if errors.Is(err, errors.New("user not found")) {
		return "", errors.New("invalid credentials")
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userQuery.Password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", errors.New("invalid credentials")
	} else if err != nil {
		return "", err
	}

	_, tokenString, _ := a.tokenAuth.Encode(Claims{"email": user.Email})

	return tokenString, nil
}

func (a *UserAuth) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.VerifyRequest(a.tokenAuth, r, jwtauth.TokenFromCookie, jwtauth.TokenFromHeader)

		if err != nil || token == nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
