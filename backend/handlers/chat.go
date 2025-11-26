package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
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

func GetChats(w http.ResponseWriter, r *http.Request) {
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

func GetChatMessages(w http.ResponseWriter, r *http.Request) {
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

func CreateChatMessage(w http.ResponseWriter, r *http.Request) {
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
}
