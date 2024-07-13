package controller

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"proxy/modules/model"
	"proxy/modules/service"
	"proxy/utils/readresponder"
)

type UserControl struct {
	readResponder readresponder.ReadResponder
	authenticator service.UserAuthenticator
}

type UserControlOption func(*UserControl)

func WithResponder(readResponder readresponder.ReadResponder) UserControlOption {
	return func(c *UserControl) {
		c.readResponder = readResponder
	}
}

func WithAuthenticator(authenticator service.UserAuthenticator) UserControlOption {
	return func(c *UserControl) {
		c.authenticator = authenticator
	}
}

func NewUserControl(options ...UserControlOption) *UserControl {
	controller := &UserControl{}

	for _, option := range options {
		option(controller)
	}

	return controller
}

// Register godoc
// @Summary register new user
// @Description Register new user provided email address and passport
// @Tags auth
// @Accept json
// @Produce json
// @Param input body User true "user credentials"
// @Success 201 {object} readresponder.JSONResponse
// @Failure 400 {object} readresponder.JSONResponse
// @Router /api/register [post]
func (c *UserControl) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := c.readResponder.ReadJSON(w, r, &user); err != nil {
		c.readResponder.WriteJSONError(w, err) // 400 status by default
	}

	if err := c.authenticator.Register(user); err != nil {
		c.readResponder.WriteJSONError(w, err) // 400 status by default
	}

	responseBody := readresponder.JSONResponse{
		Error:   false,
		Message: "user registered",
	}

	c.readResponder.WriteJSON(w, http.StatusCreated, responseBody)
}

// TODO: continue with authentication

// Authenticate godoc
// @Summary authenticate user
// @Description Authenticate user provided their email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param input body User true "user credentials"
// @Success 200 {object} JSONResponse
// @Failure 400,500 {object} JSONResponse
// @Router /api/login [post]
func (ua *UserAuth) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")

	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	err := bcrypt.CompareHashAndPassword([]byte(ua.User.Password), []byte(user.Password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || ua.User.Email != user.Email {
		w.WriteHeader(http.StatusOK)
		resp := JSONResponse{
			Code:    http.StatusOK,
			Message: "invalid credentials",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	_, tokenString, _ := ua.tokenAuth.Encode(Claims{"email": user.Email})

	w.WriteHeader(http.StatusOK)
	resp := JSONResponse{
		Code:    http.StatusOK,
		Message: "user authenticated",
		Data:    tokenString,
	}
	json.NewEncoder(w).Encode(resp)
}

func (ua *UserAuth) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwtauth.VerifyRequest(ua.tokenAuth, r, jwtauth.TokenFromCookie, jwtauth.TokenFromHeader)

		if err != nil || token == nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
