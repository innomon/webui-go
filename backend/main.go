package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/models"
	"backend/routes"
	"backend/utils"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golang-jwt/jwt/v5"
	socketio "github.com/googollee/go-socket.io"
)

// App struct represents the application
type App struct {
	Router         *chi.Mux
	SocketIOServer *socketio.Server
}

// socketIOServerAdapter adapts *socketio.Server to handlers.SocketIORoomBroadcaster
type socketIOServerAdapter struct {
	srv *socketio.Server
}

func (ad *socketIOServerAdapter) BroadcastToRoom(room string, event string, v interface{}) {
	ad.srv.BroadcastToRoom("/", room, event, v)
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
	a.SocketIOServer = socketio.NewServer(nil)

	a.SocketIOServer.OnConnect("/", func(s socketio.Conn) error {
		log.Println("Socket.IO connected:", s.ID())
		return nil
	})

	a.SocketIOServer.OnEvent("/", "auth", func(s socketio.Conn, tokenString string) {
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.GetJWTKey(), nil
		})

		if err != nil || !token.Valid {
			s.Emit("authError", "Invalid token")
			// close the connection for invalid auth
			_ = s.Close()
			return
		}

		var user models.User
		if result := database.DB.Where("email = ?", claims.Email).First(&user); result.Error != nil {
			s.Emit("authError", "User not found")
			_ = s.Close()
			return
		}

		// store the user id directly in the socket's context
		s.SetContext(user.ID)
		s.Emit("authenticated", user.ID)
		log.Printf("Socket %s authenticated for user %d\n", s.ID(), user.ID)
	})

	a.SocketIOServer.OnEvent("/", "joinChat", func(s socketio.Conn, chatIDStr string) {
		userID, ok := s.Context().(uint)
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

	a.SocketIOServer.OnEvent("/", "leaveChat", func(s socketio.Conn, chatIDStr string) {
		userID, ok := s.Context().(uint)
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

	a.SocketIOServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("Socket.IO disconnected:", s.ID(), "Reason:", reason)
	})

	a.SocketIOServer.OnError("/", func(s socketio.Conn, err error) {
		log.Println("Socket.IO error:", err)
	})

	go func() {
		if err := a.SocketIOServer.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()

	// adapter instance will be created when wiring routes
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
		// Wrap server with adapter to match handlers.SocketIORoomBroadcaster
		adapter := &socketIOServerAdapter{srv: a.SocketIOServer}
		routes.ChatRoutes(r, adapter)
		routes.LLMRoutes(r, adapter)
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
