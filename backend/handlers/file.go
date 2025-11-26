package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
)

const uploadDir = "./uploads"

func init() {
	// Create the upload directory if it doesn't exist
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	folderIDStr := r.FormValue("folder_id")
	var folderID *uint
	if folderIDStr != "" {
		id, err := strconv.ParseUint(folderIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid folder ID", http.StatusBadRequest)
			return
		}
		uID := uint(id)
		folderID = &uID

		// Verify folder belongs to user
		var folder models.Folder
		if result := database.DB.Where("id = ? AND user_id = ?", *folderID, userID).First(&folder); result.Error != nil {
			http.Error(w, "Folder not found or unauthorized", http.StatusNotFound)
			return
		}
	}

	// Create a new file on the server
	dstPath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file to server", http.StatusInternalServerError)
		return
	}

	// Save file metadata to database
	fileModel := models.File{
		UserID:   userID,
		FolderID: folderID,
		Name:     handler.Filename,
		Path:     dstPath,
		MimeType: handler.Header.Get("Content-Type"),
		Size:     handler.Size,
	}

	if result := database.DB.Create(&fileModel); result.Error != nil {
		os.Remove(dstPath) // Clean up uploaded file if DB save fails
		http.Error(w, "Failed to save file metadata", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, fileModel)
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var file models.File
	if result := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file); result.Error != nil {
		http.Error(w, "File not found or unauthorized", http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, file)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var file models.File
	if result := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file); result.Error != nil {
		http.Error(w, "File not found or unauthorized", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, file.Path)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	fileIDStr := chi.URLParam(r, "id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var file models.File
	if result := database.DB.Where("id = ? AND user_id = ?", fileID, userID).First(&file); result.Error != nil {
		http.Error(w, "File not found or unauthorized", http.StatusNotFound)
		return
	}

	// Delete from database
	if result := database.DB.Delete(&file); result.Error != nil {
		http.Error(w, "Failed to delete file from database", http.StatusInternalServerError)
		return
	}

	// Delete from file system
	if err := os.Remove(file.Path); err != nil {
		// Log the error but don't return, as the DB record is already gone
		fmt.Printf("Error deleting file from disk: %v\n", err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateFolder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var folder models.Folder
	if err := json.NewDecoder(r.Body).Decode(&folder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folder.UserID = userID

	// Check if parent folder belongs to user if specified
	if folder.ParentID != nil {
		var parentFolder models.Folder
		if result := database.DB.Where("id = ? AND user_id = ?", *folder.ParentID, userID).First(&parentFolder); result.Error != nil {
			http.Error(w, "Parent folder not found or unauthorized", http.StatusNotFound)
			return
		}
	}

	if result := database.DB.Create(&folder); result.Error != nil {
		http.Error(w, "Failed to create folder", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, folder)
}

func GetFolders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	parentIDStr := r.URL.Query().Get("parent_id")

	var folders []models.Folder
	query := database.DB.Where("user_id = ?", userID)

	if parentIDStr != "" {
		parentID, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid parent_id", http.StatusBadRequest)
			return
		}
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id IS NULL") // Top-level folders
	}

	if result := query.Find(&folders); result.Error != nil {
		http.Error(w, "Failed to retrieve folders", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, folders)
}

func GetFolderContent(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	folderIDStr := chi.URLParam(r, "id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var folder models.Folder
	if result := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder); result.Error != nil {
		http.Error(w, "Folder not found or unauthorized", http.StatusNotFound)
		return
	}

	var files []models.File
	database.DB.Where("folder_id = ? AND user_id = ?", folderID, userID).Find(&files)

	var subfolders []models.Folder
	database.DB.Where("parent_id = ? AND user_id = ?", folderID, userID).Find(&subfolders)

	response := map[string]interface{}{
		"folder":    folder,
		"files":     files,
		"subfolders": subfolders,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	folderIDStr := chi.URLParam(r, "id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid folder ID", http.StatusBadRequest)
		return
	}

	var folder models.Folder
	if result := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder); result.Error != nil {
		http.Error(w, "Folder not found or unauthorized", http.StatusNotFound)
		return
	}

	// Check for contents (files or subfolders)
	var fileCount int64
	database.DB.Model(&models.File{}).Where("folder_id = ? AND user_id = ?", folderID, userID).Count(&fileCount)
	var subfolderCount int64
	database.DB.Model(&models.Folder{}).Where("parent_id = ? AND user_id = ?", folderID, userID).Count(&subfolderCount)

	if fileCount > 0 || subfolderCount > 0 {
		http.Error(w, "Folder is not empty. Please delete all contents first.", http.StatusBadRequest)
		return
	}

	// Delete from database
	if result := database.DB.Delete(&folder); result.Error != nil {
		http.Error(w, "Failed to delete folder from database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
