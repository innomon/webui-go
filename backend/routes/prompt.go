package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// PromptRoutes defines the routes for prompt management functionality
func PromptRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Prompt routes
		r.Post("/api/prompts/create", handlers.CreatePrompt)
		r.Get("/api/prompts", handlers.GetPrompts)
		r.Get("/api/prompts/command/{command}", handlers.GetPromptByCommand)
		r.Put("/api/prompts/command/{command}/update", handlers.UpdatePrompt)
		r.Delete("/api/prompts/command/{command}/delete", handlers.DeletePrompt)
	})
}
