package dto

type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required"`
}
