package dto

// VerifyEmailRequest is the request body for verifying an email.
type VerifyAccountRequest struct {
	Token string `json:"token" validate:"required"`
}
