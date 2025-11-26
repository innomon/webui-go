package main

import (
	"log"
	"net/http"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/utils"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// App struct represents the application
type App struct {
	Router *chi.Mux
}

// Initialize initializes the application
func (a *App) Initialize() {
	database.ConnectDB()
	a.Router = chi.NewRouter()
	a.initializeMiddleware()
	a.initializeRoutes()
}

// Run starts the application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeMiddleware() {
	a.Router.Use(chi_middleware.Logger)
	a.Router.Use(chi_middleware.Recoverer)
	a.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (a *App) initializeRoutes() {
	a.Router.Get("/health", a.healthCheck)

	// Auth routes
	a.Router.Post("/api/auth/login", handlers.Login)
	a.Router.Post("/api/auth/register", handlers.Register)

	// Protected routes
	a.Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// Chat routes
		r.Post("/api/chats", handlers.CreateChat)
		r.Get("/api/chats", handlers.GetChats)
		r.Get("/api/chats/{id}/messages", handlers.GetChatMessages)
		r.Post("/api/chats/{id}/messages", handlers.CreateChatMessage)
	})
}

// handlers

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	db, err := database.DB.DB()
	if err != nil {
		http.Error(w, "Failed to get DB instance", http.StatusInternalServerError)
		return
	}

	err = db.Ping()
	if err != nil {
		http.Error(w, "DB ping failed", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func main() {
	app := App{}
	app.Initialize()
	log.Println("Starting backend server on :8080")
	app.Run(":8080")
}

