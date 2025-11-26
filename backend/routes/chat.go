package routes

import (
	"backend/handlers"
	"github.com/go-chi/chi/v5"
	socketio "github.com/doquangtan/socketio/v4"
)

// ChatRoutes defines the routes for chat functionality
func ChatRoutes(r chi.Router, srv *socketio.Server) {
	h := &handlers.Handler{SocketIOServer: srv}

	r.Post("/api/chats", h.CreateChat)
	r.Get("/api/chats", h.GetChats)
	r.Get("/api/chats/{id}/messages", h.GetChatMessages)
	r.Post("/api/chats/{id}/messages", h.CreateChatMessage)
}
