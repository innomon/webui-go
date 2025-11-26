package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

// GetModels lists all available models
func GetModels(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var models []models.Model
	query := database.DB.Where("user_id = ?", userID)

	// TODO: Implement filtering, sorting, and pagination as per python/backend/open_webui/routers/models.py

	if result := query.Find(&models); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve models"})
		return
	}

	var modelResponses []models.ModelResponse
	for _, m := range models {
		modelResponses = append(modelResponses, models.ModelResponse{
			ID:          m.ID,
			Name:        m.Name,
			BaseModelID: m.BaseModelID,
			Meta:        m.Meta,
			Params:      m.Params,
			IsActive:    m.IsActive,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ModelListResponse{
		Models: modelResponses,
		Total:  int64(len(modelResponses)), // TODO: Get actual total from filtered query
	})
}

// CreateModel creates a new model
func CreateModel(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var form models.ModelForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.ID) == "" || strings.TrimSpace(form.Name) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Model ID and Name cannot be empty"})
		return
	}

	// Check if model ID already exists
	var existingModel models.Model
	if result := database.DB.Where("id = ?", form.ID).First(&existingModel); result.RowsAffected > 0 {
		utils.RespondWithJSON(w, http.StatusConflict, map[string]string{"error": "Model with this ID already exists"})
		return
	}

	model := models.Model{
		ID:          form.ID,
		UserID:      userID,
		Name:        form.Name,
		BaseModelID: form.BaseModelID,
		Meta:        form.Meta,
		Params:      form.Params,
		AccessControl: form.AccessControl,
		IsActive:    form.IsActive,
	}

	if result := database.DB.Create(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create model"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, model)
}

// GetModelByID retrieves a single model by ID
func GetModelByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var model models.Model
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Model not found or unauthorized"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, model)
}

// UpdateModel updates an existing model
func UpdateModel(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var model models.Model
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Model not found or unauthorized"})
		return
	}

	var form models.ModelForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Name) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Model name cannot be empty"})
		return
	}

	model.Name = form.Name
	model.BaseModelID = form.BaseModelID
	model.Meta = form.Meta
	model.Params = form.Params
	model.AccessControl = form.AccessControl
	model.IsActive = form.IsActive

	if result := database.DB.Save(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update model"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, model)
}

// DeleteModel deletes a model
func DeleteModel(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")

	var model models.Model
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Model not found or unauthorized"})
		return
	}

	if result := database.DB.Delete(&model); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete model"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
