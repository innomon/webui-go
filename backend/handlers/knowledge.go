package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

// CreateKnowledge creates a new knowledge base
func CreateKnowledge(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var form models.KnowledgeForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Name) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Knowledge base name cannot be empty"})
		return
	}

	// Generate a simple collection name (slugified name)
	collectionName := strings.ToLower(strings.ReplaceAll(form.Name, " ", "-"))

	knowledge := models.Knowledge{
		UserID:       userID,
		Name:         form.Name,
		Description:  form.Description,
		CollectionName: collectionName,
		FileIDs:      []byte("[]"), // Initialize as empty JSON array
		AccessControl: []byte("{}"), // Initialize as empty JSON object
	}

	if result := database.DB.Create(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create knowledge base"})
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, knowledge)
}

// GetKnowledges retrieves all knowledge bases for a user
func GetKnowledges(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var knowledges []models.Knowledge
	if result := database.DB.Where("user_id = ?", userID).Find(&knowledges); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve knowledge bases"})
		return
	}

	// Populate files for each knowledge base (simplified, without full metadata yet)
	var knowledgeResponses []models.KnowledgeUserResponse
	for _, k := range knowledges {
		var fileIDs []uint
		json.Unmarshal(k.FileIDs, &fileIDs)

		var files []models.File
		if len(fileIDs) > 0 {
			database.DB.Where(fileIDs).Find(&files)
		}

		knowledgeResponses = append(knowledgeResponses, models.KnowledgeUserResponse{
			KnowledgeResponse: models.KnowledgeResponse{
				ID:          k.ID,
				Name:        k.Name,
				Description: k.Description,
				CreatedAt:   k.CreatedAt,
				UpdatedAt:   k.UpdatedAt,
			},
			Files: files,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, knowledgeResponses)
}

// GetKnowledgeByID retrieves a single knowledge base by ID
func GetKnowledgeByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid knowledge base ID"})
		return
	}

	var knowledge models.Knowledge
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Knowledge base not found or unauthorized"})
		return
	}

	var fileIDs []uint
	json.Unmarshal(knowledge.FileIDs, &fileIDs)

	var files []models.File
	if len(fileIDs) > 0 {
		database.DB.Where(fileIDs).Find(&files)
	}

	response := models.KnowledgeUserResponse{
		KnowledgeResponse: models.KnowledgeResponse{
			ID:          knowledge.ID,
			Name:        knowledge.Name,
			Description: knowledge.Description,
			CreatedAt:   knowledge.CreatedAt,
			UpdatedAt:   knowledge.UpdatedAt,
		},
		Files: files,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// UpdateKnowledge updates an existing knowledge base
func UpdateKnowledge(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid knowledge base ID"})
		return
	}

	var knowledge models.Knowledge
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Knowledge base not found or unauthorized"})
		return
	}

	var form models.KnowledgeForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Basic validation
	if strings.TrimSpace(form.Name) == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Knowledge base name cannot be empty"})
		return
	}

	knowledge.Name = form.Name
	knowledge.Description = form.Description

	if result := database.DB.Save(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update knowledge base"})
		return
	}

	var fileIDs []uint
	json.Unmarshal(knowledge.FileIDs, &fileIDs)

	var files []models.File
	if len(fileIDs) > 0 {
		database.DB.Where(fileIDs).Find(&files)
	}

	response := models.KnowledgeUserResponse{
		KnowledgeResponse: models.KnowledgeResponse{
			ID:          knowledge.ID,
			Name:        knowledge.Name,
			Description: knowledge.Description,
			CreatedAt:   knowledge.CreatedAt,
			UpdatedAt:   knowledge.UpdatedAt,
		},
		Files: files,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// DeleteKnowledge deletes a knowledge base
func DeleteKnowledge(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid knowledge base ID"})
		return
	}

	var knowledge models.Knowledge
	if result := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Knowledge base not found or unauthorized"})
		return
	}

	// For now, simply delete the knowledge base. In a full implementation, you would also delete
	// associated vector database collections and potentially files if not referenced elsewhere.
	if result := database.DB.Delete(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete knowledge base"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddFileToKnowledge adds a file to a knowledge base
func AddFileToKnowledge(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	knowledgeIDStr := chi.URLParam(r, "id")
	knowledgeID, err := strconv.ParseUint(knowledgeIDStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid knowledge base ID"})
		return
	}

	var knowledge models.Knowledge
	if result := database.DB.Where("id = ? AND user_id = ?", knowledgeID, userID).First(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Knowledge base not found or unauthorized"})
		return
	}

	var fileID struct { FileID uint `json:"file_id"` }
	if err := json.NewDecoder(r.Body).Decode(&fileID); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var file models.File
	if result := database.DB.Where("id = ? AND user_id = ?", fileID.FileID, userID).First(&file); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "File not found or unauthorized"})
		return
	}

	var existingFileIDs []uint
	json.Unmarshal(knowledge.FileIDs, &existingFileIDs)

	// Check if file already exists in knowledge base
	for _, id := range existingFileIDs {
		if id == fileID.FileID {
			utils.RespondWithJSON(w, http.StatusConflict, map[string]string{"error": "File already exists in knowledge base"})
			return
		}
	}

	existingFileIDs = append(existingFileIDs, fileID.FileID)
	updatedFileIDs, _ := json.Marshal(existingFileIDs)
	knowledge.FileIDs = updatedFileIDs

	if result := database.DB.Save(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add file to knowledge base"})
		return
	}

	// Re-fetch knowledge with updated files for response
	var filesInKnowledge []models.File
	database.DB.Where(existingFileIDs).Find(&filesInKnowledge)

	response := models.KnowledgeUserResponse{
		KnowledgeResponse: models.KnowledgeResponse{
			ID:          knowledge.ID,
			Name:        knowledge.Name,
			Description: knowledge.Description,
			CreatedAt:   knowledge.CreatedAt,
			UpdatedAt:   knowledge.UpdatedAt,
		},
		Files: filesInKnowledge,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

// RemoveFileFromKnowledge removes a file from a knowledge base
func RemoveFileFromKnowledge(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	knowledgeIDStr := chi.URLParam(r, "id")
	knowledgeID, err := strconv.ParseUint(knowledgeIDStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid knowledge base ID"})
		return
	}

	var knowledge models.Knowledge
	if result := database.DB.Where("id = ? AND user_id = ?", knowledgeID, userID).First(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Knowledge base not found or unauthorized"})
		return
	}

	var fileID struct { FileID uint `json:"file_id"` }
	if err := json.NewDecoder(r.Body).Decode(&fileID); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var existingFileIDs []uint
	json.Unmarshal(knowledge.FileIDs, &existingFileIDs)

	updatedFileIDs := []uint{}
	found := false
	for _, id := range existingFileIDs {
		if id == fileID.FileID {
			found = true
		} else {
			updatedFileIDs = append(updatedFileIDs, id)
		}
	}

	if !found {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "File not found in knowledge base"})
		return
	}

	updatedFileIDsBytes, _ := json.Marshal(updatedFileIDs)
	knowledge.FileIDs = updatedFileIDsBytes

	if result := database.DB.Save(&knowledge); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to remove file from knowledge base"})
		return
	}

	// Re-fetch knowledge with updated files for response
	var filesInKnowledge []models.File
	if len(updatedFileIDs) > 0 {
		database.DB.Where(updatedFileIDs).Find(&filesInKnowledge)
	}

	response := models.KnowledgeUserResponse{
		KnowledgeResponse: models.KnowledgeResponse{
			ID:          knowledge.ID,
			Name:        knowledge.Name,
			Description: knowledge.Description,
			CreatedAt:   knowledge.CreatedAt,
			UpdatedAt:   knowledge.UpdatedAt,
		},
		Files: filesInKnowledge,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}
