package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"backend/services"
	"backend/utils"
)

type SocketIORoomBroadcaster interface {
	BroadcastToRoom(room string, event string, v interface{})
}

type LLMHandler struct {
	SocketIOServer SocketIORoomBroadcaster
}

func (h *LLMHandler) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Context().Value("userID").(uint); !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var request struct {
		Model    string           `json:"model"`
		Messages []models.Message `json:"messages"`
		Stream   bool             `json:"stream"`
		ChatID   uint             `json:"chat_id"` // Added for continuity with chat history
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Determine which LLM to call based on the model name
	if strings.HasPrefix(request.Model, "ollama/") {
		allMessages := request.Messages
		if request.ChatID != 0 {
			// Fetch previous messages for context
			var previousMessages []models.Message
			database.DB.Where("chat_id = ?", request.ChatID).Order("created_at asc").Find(&previousMessages)
			allMessages = append(previousMessages, request.Messages...)
		}

		ollamaRequest := models.OllamaChatRequest{
			Model:    strings.TrimPrefix(request.Model, "ollama/"),
			Messages: allMessages,
			Stream:   request.Stream,
		}

		res, err := services.CallOllama(ollamaRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save assistant message to DB
		if request.ChatID != 0 && res != nil && res.Message.Content != "" {
			assistantMessage := models.Message{
				ChatID:  request.ChatID,
				Role:    res.Message.Role,
				Content: res.Message.Content,
			}
			if result := database.DB.Create(&assistantMessage); result.Error != nil {
				log.Printf("Error saving assistant message: %v", result.Error)
			}
			// Emit new message via Socket.IO
			h.SocketIOServer.BroadcastToRoom(fmt.Sprintf("chat:%d", request.ChatID), "message", assistantMessage)
		}

		utils.RespondWithJSON(w, http.StatusOK, res)

	} else if strings.HasPrefix(request.Model, "openai/") {
		allMessages := request.Messages
		if request.ChatID != 0 {
			// Fetch previous messages for context
			var previousMessages []models.Message
			database.DB.Where("chat_id = ?", request.ChatID).Order("created_at asc").Find(&previousMessages)
			allMessages = append(previousMessages, request.Messages...)
		}

		openaiRequest := models.OpenAIChatRequest{
			Model:    strings.TrimPrefix(request.Model, "openai/"),
			Messages: allMessages,
			Stream:   request.Stream,
		}

		res, err := services.CallOpenAI(openaiRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save assistant message to DB
		if request.ChatID != 0 && res != nil && len(res.Choices) > 0 && res.Choices[0].Message.Content != "" {
			assistantMessage := models.Message{
				ChatID:  request.ChatID,
				Role:    res.Choices[0].Message.Role,
				Content: res.Choices[0].Message.Content,
			}
			if result := database.DB.Create(&assistantMessage); result.Error != nil {
				log.Printf("Error saving assistant message: %v", result.Error)
			}
			// Emit new message via Socket.IO
			h.SocketIOServer.BroadcastToRoom(fmt.Sprintf("chat:%d", request.ChatID), "message", assistantMessage)
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	} else {
		http.Error(w, "Unsupported LLM model", http.StatusBadRequest)
		return
	}
}
