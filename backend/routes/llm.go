package routes

import (
	"backend/handlers"
	"github.com/go-chi/chi/v5"
	socketio "github.com/doquangtan/socketio/v4"
)

// LLMRoutes defines the routes for LLM functionality
func LLMRoutes(r chi.Router, srv *socketio.Server) {
	h := &handlers.LLMHandler{SocketIOServer: srv}

	r.Post("/api/chat/completions", h.ChatCompletions)
}
