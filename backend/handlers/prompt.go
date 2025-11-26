package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

// CreatePrompt creates a new prompt
func CreatePrompt(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var form models.PromptForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Command) == "" || strings.TrimSpace(form.Title) == "" || strings.TrimSpace(form.Content) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Command, Title, and Content cannot be empty"})
		return
	}

	// Check if command already exists
	var existingPrompt models.Prompt
	if result := database.DB.Where("command = ?", form.Command).First(&existingPrompt); result.RowsAffected > 0 {
		utils.RespondWithJSON(w, http.StatusConflict, map[string]string{"error": "Command with this name already exists"})
		return
	}

	prompt := models.Prompt{
		UserID:      userID,
		Title:       form.Title,
		Content:     form.Content,
		Command:     form.Command,
		AccessControl: form.AccessControl,
	}

	if result := database.DB.Create(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create prompt"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, prompt)
}

// GetPrompts retrieves all prompts for a user
func GetPrompts(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var prompts []models.Prompt
	// TODO: Implement access control logic as in python/backend/open_webui/routers/prompts.py
	if result := database.DB.Where("user_id = ?", userID).Find(&prompts); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve prompts"})
		return
	}

	var promptResponses []models.PromptUserResponse
	for _, p := range prompts {
		promptResponses = append(promptResponses, models.PromptUserResponse{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			Command:   p.Command,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, promptResponses)
}

// GetPromptByCommand retrieves a single prompt by command
func GetPromptByCommand(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	command := "/" + chi.URLParam(r, "command")

	var prompt models.Prompt
	// TODO: Implement access control logic as in python/backend/open_webui/routers/prompts.py
	if result := database.DB.Where("command = ? AND user_id = ?", command, userID).First(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Prompt not found or unauthorized"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, prompt)
}

// UpdatePrompt updates an existing prompt
func UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	command := "/" + chi.URLParam(r, "command")

	var prompt models.Prompt
	// TODO: Implement access control logic as in python/backend/open_webui/routers/prompts.py
	if result := database.DB.Where("command = ? AND user_id = ?", command, userID).First(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Prompt not found or unauthorized"})
		return
	}

	var form models.PromptForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Title) == "" || strings.TrimSpace(form.Content) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Title and Content cannot be empty"})
		return
	}

	prompt.Title = form.Title
	prompt.Content = form.Content
	prompt.AccessControl = form.AccessControl

	if result := database.DB.Save(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update prompt"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, prompt)
}

// DeletePrompt deletes a prompt
func DeletePrompt(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	command := "/" + chi.URLParam(r, "command")

	var prompt models.Prompt
	// TODO: Implement access control logic as in python/backend/open_webui/routers/prompts.py
	if result := database.DB.Where("command = ? AND user_id = ?", command, userID).First(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Prompt not found or unauthorized"})
		return
	}

	if result := database.DB.Delete(&prompt); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete prompt"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
