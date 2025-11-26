package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

// A dummy http.ResponseWriter to satisfy the interface for internal calls
type dummyResponseWriter struct{}

func (d *dummyResponseWriter) Header() http.Header        { return http.Header{} }
func (d *dummyResponseWriter) Write([]byte) (int, error)  { return 0, nil }
func (d *dummyResponseWriter) WriteHeader(statusCode int) {}

// Handler struct holds common dependencies for handlers
// Handler struct holds common dependencies for handlers
type SocketIOServer interface {
	BroadcastToRoom(room string, event string, args interface{})
}

type Handler struct {
	SocketIOServer SocketIOServer
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var chat models.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chat.UserID = userID

	if result := database.DB.Create(&chat); result.Error != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, chat)
}

func (h *Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var chats []models.Chat
	if result := database.DB.Where("user_id = ?", userID).Find(&chats); result.Error != nil {
		http.Error(w, "Failed to retrieve chats", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chats)
}

func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	chatIDStr := chi.URLParam(r, "id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var chat models.Chat
	if result := database.DB.Where("id = ? AND user_id = ?", chatID, userID).First(&chat); result.Error != nil {
		http.Error(w, "Chat not found or unauthorized", http.StatusNotFound)
		return
	}

	var messages []models.Message
	if result := database.DB.Where("chat_id = ?", chatID).Find(&messages); result.Error != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}

func (h *Handler) CreateChatMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	chatIDStr := chi.URLParam(r, "id")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var chat models.Chat
	if result := database.DB.Where("id = ? AND user_id = ?", chatID, userID).First(&chat); result.Error != nil {
		http.Error(w, "Chat not found or unauthorized", http.StatusNotFound)
		return
	}

	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message.ChatID = uint(chatID)

	if result := database.DB.Create(&message); result.Error != nil {
		http.Error(w, "Failed to create message", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, message)

	// Broadcast the new message via Socket.IO
	h.SocketIOServer.BroadcastToRoom(fmt.Sprintf("chat:%d", chatID), "message", message)

	// After saving the user's message, call the LLM to get a response
	var allMessages []models.Message
	database.DB.Where("chat_id = ?", chatID).Order("created_at asc").Find(&allMessages)

	llmRequest := struct {
		Model    string           `json:"model"`
		Messages []models.Message `json:"messages"`
		Stream   bool             `json:"stream"`
		ChatID   uint             `json:"chat_id"`
	}{
		Model:    "ollama/llama3", // Defaulting to an Ollama model for now
		Messages: allMessages,
		Stream:   false,
		ChatID:   uint(chatID),
	}

	requestBody, err := json.Marshal(llmRequest)
	if err != nil {
		log.Printf("Error marshaling LLM request: %v", err)
		return
	}

	// Create a new HTTP request to the LLM completions endpoint
	req, err := http.NewRequest("POST", "/api/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating LLM request: %v", err)
		return
	}

	// Set context with userID for authentication in the LLM handler
	req = req.WithContext(r.Context())

	// Create a new ResponseRecorder to capture the LLM handler's response
	rr := httptest.NewRecorder()

	// Create an LLMHandler instance and call its ChatCompletions method
	llmHandler := LLMHandler{SocketIOServer: h.SocketIOServer}
	llmHandler.ChatCompletions(rr, req)

	// The LLM handler will save the message and broadcast it via Socket.IO
}
