package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/utils"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	socketio "github.com/doquangtan/socketio/v4"
	"github.com/golang-jwt/jwt/v5"
)

// App struct represents the application
type App struct {
	Router        *chi.Mux
	SocketIOServer *socketio.Server
}

// Initialize initializes the application
func (a *App) Initialize() {
	database.ConnectDB()
	a.Router = chi.NewRouter()
	a.initializeMiddleware()
	a.initializeSocketIO()
	a.initializeRoutes()
}

// Run starts the application
func (a *App) Run(addr string) {
	defer a.SocketIOServer.Close()
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

func (a *App) initializeSocketIO() {
	var err error
	a.SocketIOServer = socketio.NewServer(nil)
	if err != nil {
		log.Fatalf("Failed to create Socket.IO server: %v", err)
	}

	a.SocketIOServer.OnConnect(func(s socketio.Socket) error {
		log.Println("Socket.IO connected:", s.ID())
		return nil
	})

	a.SocketIOServer.OnEvent("auth", func(s socketio.Socket, tokenString string) {
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.GetJWTKey(), nil
		})

		if err != nil || !token.Valid {
			s.Emit("authError", "Invalid token")
			s.Disconnect()
			return
		}

		var user models.User
		if result := database.DB.Where("email = ?", claims.Email).First(&user); result.Error != nil {
			s.Emit("authError", "User not found")
			s.Disconnect()
			return
		}

		s.SetContext(context.WithValue(s.Context(), "userID", user.ID))
		s.Emit("authenticated", user.ID)
		log.Printf("Socket %s authenticated for user %d\n", s.ID(), user.ID)
	})

	a.SocketIOServer.OnEvent("joinChat", func(s socketio.Socket, chatIDStr string) {
		userID, ok := s.Context().Value("userID").(uint)
		if !ok {
			s.Emit("error", "Unauthorized")
			return
		}

		chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
		if err != nil {
			s.Emit("error", "Invalid chat ID")
			return
		}

		var chat models.Chat
		if result := database.DB.Where("id = ? AND user_id = ?", chatID, userID).First(&chat); result.Error != nil {
			s.Emit("error", "Chat not found or unauthorized")
			return
		}

		s.Join(fmt.Sprintf("chat:%d", chatID))
		log.Printf("Socket %s joined chat %d for user %d\n", s.ID(), chatID, userID)
		s.Emit("joinedChat", chatID)
	})

	a.SocketIOServer.OnEvent("leaveChat", func(s socketio.Socket, chatIDStr string) {
		userID, ok := s.Context().Value("userID").(uint)
		if !ok {
			s.Emit("error", "Unauthorized")
			return
		}
		chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
		if err != nil {
			s.Emit("error", "Invalid chat ID")
			return
		}
		s.Leave(fmt.Sprintf("chat:%d", chatID))
		log.Printf("Socket %s left chat %d for user %d\n", s.ID(), chatID, userID)
		s.Emit("leftChat", chatID)
	})

	a.SocketIOServer.OnDisconnect(func(s socketio.Socket, reason string) {
		log.Println("Socket.IO disconnected:", s.ID(), "Reason:", reason)
	})

	a.SocketIOServer.OnError(func(s socketio.Socket, err error) {
		log.Println("Socket.IO error:", err)
	})

	go func() {
		if err := a.SocketIOServer.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()
}

func (a *App) initializeRoutes() {
	a.Router.Get("/health", a.healthCheck)

	// Auth routes
	a.Router.Post("/api/auth/login", handlers.Login)
	a.Router.Post("/api/auth/register", handlers.Register)

	// Mount Socket.IO server
	a.Router.Handle("/socket.io/*", a.SocketIOServer)

	// Protected routes
	a.Router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		routes.ChatRoutes(r, a.SocketIOServer)
		routes.LLMRoutes(r, a.SocketIOServer)
		routes.FileRoutes(r)
		routes.KnowledgeRoutes(r)
		routes.ModelRoutes(r)
		routes.PromptRoutes(r)
		routes.ToolRoutes(r)
		routes.UserAdminRoutes(r)
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

