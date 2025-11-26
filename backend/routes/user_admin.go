package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/go-chi/chi/v5"
)

// UserAdminRoutes defines the routes for user administration functionality
func UserAdminRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		// These routes should be protected and only accessible by admins.
		// For simplicity, we use the general AuthMiddleware. A more robust solution
		// would involve a specific AdminMiddleware.
		r.Use(middleware.AuthMiddleware)

		// User administration routes
		r.Get("/api/users", handlers.GetUsers)
		r.Put("/api/users/{id}", handlers.UpdateUser)
		r.Delete("/api/users/{id}", handlers.DeleteUser)

		// Route to get the current user's profile
		r.Get("/api/user/me", handlers.GetCurrentUser)
	})
}
