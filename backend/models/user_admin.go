package models

// UserUpdateForm for updating a user's profile
type UserUpdateForm struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password,omitempty"`
	ProfileImageURL string `json:"profile_image_url"`
	Role            string `json:"role"`
}

// UserRoleUpdateForm for updating a user's role
type UserRoleUpdateForm struct {
	Role string `json:"role" binding:"required"`
}

// UserSettings represents user-specific settings
type UserSettings struct {
	UI map[string]interface{} `json:"ui"`
}

// UserGroupIdsModel for returning user info with group IDs
type UserGroupIdsModel struct {
	User
	GroupIDs []uint `json:"group_ids"`
}

// UserGroupIdsListResponse for returning a list of users with group IDs and total count
type UserGroupIdsListResponse struct {
	Users []UserGroupIdsModel `json:"users"`
	Total int64             `json:"total"`
}

// UserInfoListResponse for returning a list of users with total count
type UserInfoListResponse struct {
	Users []User `json:"users"`
	Total int64  `json:"total"`
}

// UserIdNameListResponse for returning a list of user IDs and names with total count
type UserIdNameListResponse struct {
	Users []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"users"`
	Total int64 `json:"total"`
}
