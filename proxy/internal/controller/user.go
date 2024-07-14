package controller

import (
	"net/http"
	"proxy/internal/entities"
	"proxy/internal/service"
	"proxy/utils/readresponder"
)

type UserControl struct {
	readResponder readresponder.ReadResponder
	authService   service.Authenticator
}

type UserControlOption func(*UserControl)

func WithResponder(readResponder readresponder.ReadResponder) UserControlOption {
	return func(c *UserControl) {
		c.readResponder = readResponder
	}
}

func WithAuthenticator(authService service.Authenticator) UserControlOption {
	return func(c *UserControl) {
		c.authService = authService
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
// @Param input body entities.User true "user credentials"
// @Success 201 {object} readresponder.JSONResponse
// @Failure 400 {object} readresponder.JSONResponse
// @Router /api/register [post]
func (c *UserControl) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	if err := c.readResponder.ReadJSON(w, r, &user); err != nil {
		c.readResponder.WriteJSONError(w, err) // 400 status by default
		return
	}

	if err := c.authService.Register(user); err != nil {
		c.readResponder.WriteJSONError(w, err) // 400 status by default
		return
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
// @Param input body entities.User true "user credentials"
// @Success 200 {object} readresponder.JSONResponse
// @Failure 400,500 {object} readresponder.JSONResponse
// @Router /api/login [post]
func (c *UserControl) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := c.readResponder.ReadJSON(w, r, &user); err != nil {
		c.readResponder.WriteJSONError(w, err)
		return
	}

	tokenString, err := c.authService.Authenticate(user)
	if err != nil {
		c.readResponder.WriteJSONError(w, err)
		return
	}

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "user authenticated",
		Data:    tokenString,
	}

	c.readResponder.WriteJSON(w, http.StatusOK, resp)
}

func (c *UserControl) RequireAuthentication(next http.Handler) http.Handler {
	return c.authService.RequireAuthentication(next)
}
