package services

import (
	"context"
	dto "scs-user/internal/dto"
	repositories "scs-user/internal/repositories"
	"scs-user/pkg/errors"
	"scs-user/pkg/utils"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, loginDto *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, loginDto.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("User not found")
	}

	// Verify the password
	if err := utils.VerifyPassword(user.Password, loginDto.Password); err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.NewUnauthorizedError("User is not active")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate token", err)
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	err := utils.ValidateToken(token)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid token")
	}
	return &dto.ValidateTokenResponse{
		Valid: true,
	}, nil
}
