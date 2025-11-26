package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// FileRoutes defines the routes for file and folder functionality
func FileRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// File routes
		r.Post("/api/files/upload", handlers.UploadFile)
		r.Get("/api/files/{id}", handlers.GetFile)
		r.Get("/api/files/{id}/download", handlers.DownloadFile)
		r.Delete("/api/files/{id}", handlers.DeleteFile)

		// Folder routes
		r.Post("/api/folders", handlers.CreateFolder)
		r.Get("/api/folders", handlers.GetFolders)
		r.Get("/api/folders/{id}", handlers.GetFolderContent)
		r.Delete("/api/folders/{id}", handlers.DeleteFolder)
	})
}
