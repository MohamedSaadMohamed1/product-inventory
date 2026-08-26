package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	sharedHttp "product-inventory/internal/http"
)

// Handler handles HTTP requests for products.
type Handler struct {
	service *Service
}

// NewHandler creates a new Handler instance.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles the product creation request.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body")
		return
	}

	p, err := h.service.Create(r.Context(), req)
	if err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	sharedHttp.WriteJSON(w, http.StatusCreated, p)
}

// Get handles retrieving a single product by ID.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	p, err := h.service.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			sharedHttp.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, p)
}

// List handles paginated product listing.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := int32(1)
	limit := int32(20)

	if pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 32)
		if err == nil && p > 0 {
			page = int32(p)
		}
	}
	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err == nil && l > 0 {
			limit = int32(l)
		}
	}

	products, err := h.service.List(r.Context(), page, limit)
	if err != nil {
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, products)
}

// Update handles product updates.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to decode request body")
		return
	}

	p, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			sharedHttp.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		sharedHttp.WriteError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	sharedHttp.WriteJSON(w, http.StatusOK, p)
}

// Delete handles deleting a product.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		sharedHttp.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			sharedHttp.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		sharedHttp.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
