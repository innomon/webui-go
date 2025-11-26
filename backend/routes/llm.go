package routes

import (
	"backend/handlers"

	"github.com/go-chi/chi/v5"
)

// LLMRoutes defines the routes for LLM functionality
func LLMRoutes(r chi.Router, srv handlers.SocketIORoomBroadcaster) {
	h := &handlers.LLMHandler{SocketIOServer: srv}

	r.Post("/api/chat/completions", h.ChatCompletions)
}
