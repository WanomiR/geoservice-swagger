package controller

import (
	"net/http"
	"proxy/internal/modules/auth/entities"
	"proxy/internal/modules/auth/service"
	"proxy/internal/utils/readresponder"
)

type Auth struct {
	authService   service.Authenticator
	readResponder readresponder.ReadResponder
}

func NewAuth(authService service.Authenticator, responder readresponder.ReadResponder) *Auth {
	return &Auth{authService: authService, readResponder: responder}
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
func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	if err := a.readResponder.ReadJSON(w, r, &user); err != nil {
		a.readResponder.WriteJSONError(w, err) // 400 status by default
		return
	}

	if user.Password == "" || user.Email == "" {
		a.readResponder.WriteJSONError(w, service.ErrorEOF)
		return
	}

	if err := a.authService.Register(user); err != nil {
		a.readResponder.WriteJSONError(w, err) // 400 status by default
		return
	}

	responseBody := readresponder.JSONResponse{
		Error:   false,
		Message: "user registered",
	}

	a.readResponder.WriteJSON(w, http.StatusCreated, responseBody)
}

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
func (a *Auth) Authenticate(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	if err := a.readResponder.ReadJSON(w, r, &user); err != nil {
		a.readResponder.WriteJSONError(w, err)
		return
	}

	if user.Password == "" || user.Email == "" {
		a.readResponder.WriteJSONError(w, service.ErrorEOF)
		return
	}

	tokenString, err := a.authService.Authenticate(user)
	if err != nil {
		a.readResponder.WriteJSONError(w, err)
		return
	}

	resp := readresponder.JSONResponse{
		Error:   false,
		Message: "user authenticated",
		Data:    tokenString,
	}

	a.readResponder.WriteJSON(w, http.StatusOK, resp)
}
