package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// App struct represents the application
type App struct {
	Router *http.ServeMux
}

// Initialize initializes the application
func (a *App) Initialize() {
	a.Router = http.NewServeMux()
	a.initializeRoutes()
}

// Run starts the application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	// Auth routes
	a.Router.HandleFunc("/api/auth/login", a.login)
	a.Router.HandleFunc("/api/auth/register", a.register)

	// Chat routes
	a.Router.HandleFunc("/api/chats", a.chatsHandler)
	a.Router.HandleFunc("/api/chats/", a.chatHandler) // Note the trailing slash
}

// handlers

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "logged in"})
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "registered"})
}

func (a *App) chatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/chats" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		a.getChats(w, r)
	case http.MethodPost:
		a.createChat(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) chatHandler(w http.ResponseWriter, r *http.Request) {
	// This will handle /api/chats/{id} and /api/chats/{id}/messages
	path := strings.TrimPrefix(r.URL.Path, "/api/chats/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	chatID := parts[0]

	if len(parts) == 1 { // Route: /api/chats/{id}
		if r.Method == http.MethodGet {
			a.getChat(w, r, chatID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if len(parts) == 2 && parts[1] == "messages" { // Route: /api/chats/{id}/messages
		switch r.Method {
		case http.MethodGet:
			a.getChatMessages(w, r, chatID)
		case http.MethodPost:
			a.createChatMessage(w, r, chatID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.NotFound(w, r)
	}
}

func (a *App) getChats(w http.ResponseWriter, r *http.Request) {
	// Placeholder data
	chats := []map[string]interface{}{
		{"id": "1", "name": "Chat 1"},
		{"id": "2", "name": "Chat 2"},
	}
	respondWithJSON(w, http.StatusOK, chats)
}

func (a *App) createChat(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "chat created"})
}

func (a *App) getChat(w http.ResponseWriter, r *http.Request, chatID string) {
	// Placeholder data
	chat := map[string]interface{}{"id": chatID, "name": "Chat " + chatID}
	respondWithJSON(w, http.StatusOK, chat)
}

func (a *App) getChatMessages(w http.ResponseWriter, r *http.Request, chatID string) {
	// Placeholder data
	messages := []map[string]interface{}{
		{"id": "1", "chat_id": chatID, "text": "Hello"},
		{"id": "2", "chat_id": chatID, "text": "World"},
	}
	respondWithJSON(w, http.StatusOK, messages)
}

func (a *App) createChatMessage(w http.ResponseWriter, r *http.Request, chatID string) {
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "message created"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	app := App{}
	app.Initialize()
	log.Println("Starting backend server on :8080")
	app.Run(":8080")
}