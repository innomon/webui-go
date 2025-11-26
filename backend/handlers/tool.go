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

// CreateTool creates a new tool
func CreateTool(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var form models.ToolForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.ID) == "" || strings.TrimSpace(form.Name) == "" || strings.TrimSpace(form.Content) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "ID, Name, and Content cannot be empty"})
		return
	}

	// Check if tool ID already exists
	var existingTool models.Tool
	if result := database.DB.Where("id = ?", form.ID).First(&existingTool); result.RowsAffected > 0 {
		utils.RespondWithJSON(w, http.StatusConflict, map[string]string{"error": "Tool with this ID already exists"})
		return
	}

	// TODO: Implement Python code processing and spec generation as in python/backend/open_webui/routers/tools.py

	tool := models.Tool{
		ID:        form.ID,
		UserID:    userID,
		Name:      form.Name,
		Content:   form.Content,
		Meta:      form.Meta,
		AccessControl: form.AccessControl,
		Specs:     []byte("{}"), // Placeholder for generated specs
	}

	if result := database.DB.Create(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create tool"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, tool)
}

// GetTools retrieves all tools for a user
func GetTools(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var tools []models.Tool
	// TODO: Implement access control logic as in python/backend/open_webui/routers/tools.py
	if result := database.DB.Where("user_id = ?", userID).Find(&tools); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve tools"})
		return
	}

	var toolResponses []models.ToolUserResponse
	for _, t := range tools {
		toolResponses = append(toolResponses, models.ToolUserResponse{
			ToolResponse: models.ToolResponse{
				ID:        t.ID,
				Name:      t.Name,
				Content:   t.Content,
				Specs:     t.Specs,
				Meta:      t.Meta,
				CreatedAt: t.CreatedAt,
				UpdatedAt: t.UpdatedAt,
			},
			// TODO: Implement has_user_valves logic
			HasUserValves: false,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, toolResponses)
}

// GetToolByID retrieves a single tool by ID
func GetToolByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var tool models.Tool
	// TODO: Implement access control logic as in python/backend/open_webui/routers/tools.py
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Tool not found or unauthorized"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tool)
}

// UpdateTool updates an existing tool
func UpdateTool(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var tool models.Tool
	// TODO: Implement access control logic as in python/backend/open_webui/routers/tools.py
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Tool not found or unauthorized"})
		return
	}

	var form models.ToolForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Name) == "" || strings.TrimSpace(form.Content) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Name and Content cannot be empty"})
		return
	}

	// TODO: Implement Python code processing and spec generation as in python/backend/open_webui/routers/tools.py

	tool.Name = form.Name
	tool.Content = form.Content
	tool.Meta = form.Meta
	tool.AccessControl = form.AccessControl

	if result := database.DB.Save(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update tool"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, tool)
}

// DeleteTool deletes a tool
func DeleteTool(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var tool models.Tool
	// TODO: Implement access control logic as in python/backend/open_webui/routers/tools.py
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Tool not found or unauthorized"})
		return
	}

	if result := database.DB.Delete(&tool); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete tool"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
