package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// ToolRoutes defines the routes for tool management functionality
func ToolRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Tool routes
		r.Post("/api/tools/create", handlers.CreateTool)
		r.Get("/api/tools", handlers.GetTools)
		r.Get("/api/tools/id/{id}", handlers.GetToolByID)
		r.Put("/api/tools/id/{id}/update", handlers.UpdateTool)
		r.Delete("/api/tools/id/{id}/delete", handlers.DeleteTool)

		// TODO: Implement other tool-related endpoints (valves, user valves)
	})
}
