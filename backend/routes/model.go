package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// ModelRoutes defines the routes for model management functionality
func ModelRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Model routes
		r.Post("/api/models/create", handlers.CreateModel)
		r.Get("/api/models/list", handlers.GetModels)
		r.Get("/api/models/{id}", handlers.GetModelByID)
		r.Put("/api/models/{id}", handlers.UpdateModel)
		r.Delete("/api/models/{id}", handlers.DeleteModel)

		// TODO: Implement other model-related endpoints (base models, tags, export, import, sync, profile image)
	})
}
