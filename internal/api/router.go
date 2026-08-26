package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"product-inventory/internal/auth"
	authmiddleware "product-inventory/internal/middleware"
	"product-inventory/internal/order"
	"product-inventory/internal/product"
)

// NewRouter registers application routes, mounts auth handlers, and secures endpoints using JWT middleware.
func NewRouter(
	productHandler *product.Handler,
	orderHandler *order.Handler,
	authHandler *auth.Handler,
	jwtSecret string,
) *chi.Mux {
	r := chi.NewRouter()

	// General middleware stack
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(authmiddleware.RateLimit(10, 20)) // Allow up to 10 requests/sec with burst capacity of 20

	// Static routes
	r.Get("/", ServeDashboard)
	r.Get("/swagger", ServeSwaggerUI)
	r.Get("/docs/openapi.yaml", ServeOpenAPIYAML)

	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes (Public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// Product routes
		r.Route("/products", func(r chi.Router) {
			// Public reads
			r.Get("/", productHandler.List)
			r.Get("/{id}", productHandler.Get)

			// Admin-only writes
			r.Group(func(r chi.Router) {
				r.Use(authmiddleware.Authenticate(jwtSecret))
				r.Use(authmiddleware.RequireRole("admin"))

				r.Post("/", productHandler.Create)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})
		})

		// Order routes (All require authentication)
		r.Route("/orders", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authmiddleware.Authenticate(jwtSecret))

				r.Get("/", orderHandler.List)
				r.Post("/", orderHandler.Create)
				r.Get("/{id}", orderHandler.Get)
				r.Post("/{id}/cancel", orderHandler.Cancel)
			})
		})
	})

	return r
}
