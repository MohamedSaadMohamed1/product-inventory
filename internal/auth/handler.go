package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	sharedHttp "product-inventory/internal/http"
)

// Handler handles HTTP requests for user authentication.
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration request.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body")
		return
	}

	res, err := h.service.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUsernameTaken) {
			sharedHttp.WriteError(w, http.StatusConflict, "USERNAME_TAKEN", err.Error())
			return
		}
		sharedHttp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	sharedHttp.WriteJSON(w, http.StatusCreated, res)
}

// Login handles user authentication request.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body")
		return
	}

	res, err := h.service.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			sharedHttp.WriteError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
			return
		}
		sharedHttp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, res)
}
