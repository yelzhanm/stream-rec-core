package http

import (
	"encoding/json"
	"net/http"

	"stream-rec-core/internal/domain"
	"stream-rec-core/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	svc *usecase.PlatformService
	log zerolog.Logger
}

func NewHandler(svc *usecase.PlatformService, log zerolog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

func (h *Handler) respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// getProfile godoc
// @Summary      Get user profile
// @Description  Retrieves user profile details by MongoDB ObjectID
// @Tags         Users & Auth
// @Produce      json
// @Param        id   path      string  true  "User ObjectID"
// @Success      200  {object}  domain.User
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	u, err := h.svc.GetProfile(r.Context(), id)
	if err != nil {
		h.respond(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	h.respond(w, http.StatusOK, u)
}

// recordInteraction godoc
// @Summary      Record user interaction
// @Description  Logs a user interaction with an item (requires JWT Bearer Token)
// @Tags         Interactions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        interaction  body      domain.Interaction  true  "Interaction event payload"
// @Success      201          {object}  domain.Interaction
// @Failure      401          {object}  ErrorResponse
// @Failure      500          {object}  ErrorResponse
// @Router       /interactions [post]
func (h *Handler) recordInteraction(w http.ResponseWriter, r *http.Request) {
	var inter domain.Interaction
	if err := json.NewDecoder(r.Body).Decode(&inter); err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": "invalid payload"})
		return
	}

	if claims, ok := r.Context().Value(userContextKey).(*domain.JWTClaims); ok && inter.UserID.IsZero() {
		if uid, err := primitive.ObjectIDFromHex(claims.UserID); err == nil {
			inter.UserID = uid
		}
	}

	if err := h.svc.RecordInteraction(r.Context(), &inter); err != nil {
		h.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusCreated, inter)
}

// getUserHistory godoc
// @Summary      Get user interaction history
// @Description  Retrieves interaction history for a user (requires JWT Bearer Token)
// @Tags         Interactions
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      string  true  "User ObjectID"
// @Success      200  {array}   domain.Interaction
// @Failure      401  {object}  ErrorResponse
// @Router       /users/{id}/interactions [get]
func (h *Handler) getUserHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	history, err := h.svc.GetUserHistory(r.Context(), id)
	if err != nil {
		h.respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	h.respond(w, http.StatusOK, history)
}
