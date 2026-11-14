package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"stream-rec-core/internal/domain"
)

type contextKey string

const userContextKey contextKey = "authUser"

// register godoc
// @Summary      Register a new user
// @Description  Creates a new user account and profile with favorite genres
// @Tags         Users & Auth
// @Accept       json
// @Produce      json
// @Param        user  body      domain.User  true  "User registration payload"
// @Success      201   {object}  domain.User
// @Failure      400   {object}  ErrorResponse
// @Failure      409   {object}  ErrorResponse
// @Router       /users/register [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var u domain.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}
	if err := h.svc.Register(r.Context(), &u); err != nil {
		h.log.Error().Err(err).Msg("failed to register user")
		h.respond(w, http.StatusConflict, map[string]string{"error": "email already exists"})
		return
	}
	h.respond(w, http.StatusCreated, u)
}

type LoginRequest struct {
	Email    string `json:"email" example:"yelzhan@aitu.edu.kz"`
	Password string `json:"password" example:"secret123"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// login godoc
// @Summary      Authenticate user
// @Description  Logs in a user and returns a signed JWT token
// @Tags         Users & Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Login credentials"
// @Success      200          {object}  LoginResponse
// @Failure      401          {object}  ErrorResponse
// @Router       /users/login [post]
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&creds)
	token, err := h.svc.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		h.respond(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, LoginResponse{
		Token: token,
	})
}

func (h *Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.respond(w, http.StatusUnauthorized, ErrorResponse{Error: "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			h.respond(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid authorization format, expected: Bearer <token>"})
			return
		}

		claims, err := h.svc.ValidateToken(parts[1])
		if err != nil {
			h.respond(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized: " + err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
