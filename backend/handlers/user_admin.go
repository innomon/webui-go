package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backend/database"
	"backend/models"
	"backend/utils"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers retrieves a paginated list of users (admin only)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// For simplicity, this is a basic implementation.
	// TODO: Add filtering, sorting, and pagination as in the Python backend.

	var users []models.User
	if result := database.DB.Find(&users); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve users"})
		return
	}

	// For now, returning a simplified list without group IDs.
	utils.RespondWithJSON(w, http.StatusOK, models.UserInfoListResponse{
		Users: users,
		Total: int64(len(users)),
	})
}

// UpdateUser updates a user's information (admin only)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	var form models.UserUpdateForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user.Name = form.Name
	user.Email = form.Email
	user.Role = form.Role
	// user.ProfileImageURL = form.ProfileImageURL // Add this field to the User model if needed

	if form.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), 8)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if result := database.DB.Save(&user); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser deletes a user (admin only)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	// It's good practice to prevent self-deletion or deletion of a primary admin.
	// TODO: Add logic to prevent deletion of the primary admin user.

	if result := database.DB.Delete(&models.User{}, userID); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetCurrentUser retrieves the profile of the currently authenticated user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		utils.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	// Return a subset of user info, excluding sensitive data like password
	userInfo := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		// "profile_image_url": user.ProfileImageURL,
	}

	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}
