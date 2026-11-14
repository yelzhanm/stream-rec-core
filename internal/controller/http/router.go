package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (h *Handler) InitRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", h.register)
			r.Post("/login", h.login)
			r.Get("/{id}", h.getProfile)
		})

		r.Route("/items", func(r chi.Router) {
			r.Get("/search", h.searchItems)
			r.Post("/", h.createItem)
			r.Get("/", h.getItems)
			r.Get("/{id}", h.getItemByID)
			r.Put("/{id}", h.updateItem)
			r.Delete("/{id}", h.deleteItem)
		})

		r.Group(func(protected chi.Router) {
			protected.Use(h.authMiddleware)

			protected.Post("/interactions", h.recordInteraction)
			protected.Get("/users/{id}/interactions", h.getUserHistory)
		})
	})

	return r
}
