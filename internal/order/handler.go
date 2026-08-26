package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	sharedHttp "product-inventory/internal/http"
	"product-inventory/internal/middleware"
)

// Handler handles HTTP requests for orders.
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles requests to create a new order and reserve stock.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body")
		return
	}

	order, err := h.service.Create(r.Context(), claims.UserID, req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			sharedHttp.WriteError(w, http.StatusBadRequest, "PRODUCT_NOT_FOUND", err.Error())
			return
		}
		if errors.Is(err, ErrInsufficientStock) {
			sharedHttp.WriteError(w, http.StatusConflict, "INSUFFICIENT_STOCK", err.Error())
			return
		}
		sharedHttp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	sharedHttp.WriteJSON(w, http.StatusCreated, order)
}

// Get handles requests to fetch an order by ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	order, err := h.service.Get(r.Context(), id, claims.UserID, claims.Role)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			sharedHttp.WriteError(w, http.StatusNotFound, "ORDER_NOT_FOUND", "Order not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			sharedHttp.WriteError(w, http.StatusForbidden, "FORBIDDEN_ACCESS", "You do not have permission to view this order")
			return
		}
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, order)
}

// Cancel handles requests to cancel a pending order.
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid order ID")
		return
	}

	err = h.service.Cancel(r.Context(), id, claims.UserID, claims.Role)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			sharedHttp.WriteError(w, http.StatusNotFound, "ORDER_NOT_FOUND", "Order not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			sharedHttp.WriteError(w, http.StatusForbidden, "FORBIDDEN_ACCESS", "You do not have permission to cancel this order")
			return
		}
		if errors.Is(err, ErrOrderAlreadyCancelled) {
			sharedHttp.WriteError(w, http.StatusConflict, "ORDER_ALREADY_CANCELLED", "Order is already cancelled")
			return
		}
		if errors.Is(err, ErrInvalidOrderState) {
			sharedHttp.WriteError(w, http.StatusConflict, "INVALID_ORDER_STATE", err.Error())
			return
		}
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Order cancelled successfully",
	})
}

// List handles requests to list the authenticated user's orders with pagination.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		sharedHttp.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication required")
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	orders, err := h.service.List(r.Context(), claims.UserID, page, limit)
	if err != nil {
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve orders")
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, orders)
}
