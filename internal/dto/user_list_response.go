package dto

import (
	"scs-user/internal/models"
	"scs-user/internal/types"
)

// UserListResponse represents the paginated response for user list
type UserListResponse struct {
	Data       []models.User     `json:"data"`
	Pagination types.Pagination `json:"pagination"`
}
