package auth

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type User struct {
	Email    string `json:"email" binding:"required" example:"admin@example.com"`
	Password string `json:"password" binding:"required" example:"password"`
}

type Claims map[string]interface{}

type UserAuth struct {
	User      User
	tokenAuth *jwtauth.JWTAuth
}

type JSONResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
	Data    any    `json:"data,omitempty"`
}

func NewUserAuth(algorithm, secret string) *UserAuth {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	tokenAuth := jwtauth.New(algorithm, []byte(secret), nil)

	userAuth := &UserAuth{
		User{"admin@example.com", string(encryptedPassword)},
		tokenAuth,
	}
	return userAuth
}

// RegisterUser godoc
// @Summary register user
// @Description Register new user provided email address and passport
// @Tags auth
// @Accept json
// @Produce json
// @Param input body User true "user credentials"
// @Success 201 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /api/register [post]
func (ua *UserAuth) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp := JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer r.Body.Close()

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	ua.User.Password = string(encryptedPassword)
	ua.User.Email = user.Email

	resp := JSONResponse{
		Code:    http.StatusCreated,
		Message: "user registered",
	}
	json.NewEncoder(w).Encode(resp)
}

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

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		resp := JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	defer r.Body.Close()

	err = bcrypt.CompareHashAndPassword([]byte(ua.User.Password), []byte(user.Password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || ua.User.Email != user.Email {
		w.WriteHeader(http.StatusOK)
		resp := JSONResponse{
			Code:    http.StatusOK,
			Message: "invalid credentials",
		}
		json.NewEncoder(w).Encode(resp)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
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
