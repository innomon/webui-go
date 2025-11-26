package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// KnowledgeRoutes defines the routes for knowledge base functionality
func KnowledgeRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Knowledge base routes
		r.Post("/api/knowledge/create", handlers.CreateKnowledge)
		r.Get("/api/knowledge", handlers.GetKnowledges)
		r.Get("/api/knowledge/{id}", handlers.GetKnowledgeByID)
		r.Put("/api/knowledge/{id}", handlers.UpdateKnowledge)
		r.Delete("/api/knowledge/{id}", handlers.DeleteKnowledge)
		r.Post("/api/knowledge/{id}/file/add", handlers.AddFileToKnowledge)
		r.Post("/api/knowledge/{id}/file/remove", handlers.RemoveFileFromKnowledge)
	})
}
