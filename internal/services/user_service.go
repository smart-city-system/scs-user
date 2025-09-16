package services

import (
	"context"
	"encoding/json"
	dto "scs-user/internal/dto"
	"scs-user/internal/models"
	repositories "scs-user/internal/repositories"
	"scs-user/internal/types"
	"scs-user/pkg/errors"
	kafka_client "scs-user/pkg/kafka"
	"scs-user/pkg/utils"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type UserService struct {
	userRepo        repositories.UserRepository
	userPremiseRepo repositories.UserPremiseRepository
	producer        kafka_client.Producer
}

func NewUserService(userRepo repositories.UserRepository, userPremiseRepo repositories.UserPremiseRepository, producer kafka_client.Producer) *UserService {
	return &UserService{userRepo: userRepo, userPremiseRepo: userPremiseRepo, producer: producer}
}

func (s *UserService) CreateUser(ctx context.Context, createUserDto *dto.CreateUserDto) (*models.User, error) {
	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(createUserDto.Password)
	if err != nil {
		return nil, errors.NewInternalError("Failed to hash password", err)
	}

	user := &models.User{
		Name:     createUserDto.Name,
		Email:    createUserDto.Email,
		Password: hashedPassword,
		Role:     createUserDto.Role,
		IsActive: false,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		// Check if it's a duplicate email error
		if isDuplicateEmailError(err) {
			return nil, errors.NewConflictError("User with this email already exists")
		}
		return nil, errors.NewDatabaseError("create user", err)
	}
	if createUserDto.PremiseID != "" {
		// Assign the user to the premise
		premiseId, err := uuid.Parse(createUserDto.PremiseID)
		if err != nil {
			return nil, errors.NewBadRequestError("Invalid premise id")
		}
		userPremise := &models.UserPremise{
			UserID:    createdUser.ID,
			PremiseID: premiseId,
		}
		err = s.userPremiseRepo.AssignPremises(ctx, userPremise)
		if err != nil {
			return nil, errors.NewDatabaseError("add user to premise", err)
		}
	}
	// Generate JWT token
	token, err := utils.GenerateToken(createdUser.ID.String(), createdUser.Role)
	if err != nil {
		return nil, errors.NewInternalError("Failed to generate token", err)
	}
	//send to kafka
	message := types.Message[map[string]interface{}]{
		Type:    "user.created",
		Payload: map[string]interface{}{"token": token, "email": createdUser.Email},
	}
	messageBytes, err := json.Marshal(message)

	// Send Kafka message after successful alarm creation

	if err != nil {
		// Log error but don't fail the operation since alarm was created successfully
		// You might want to add proper logging here
		return nil, errors.NewInternalError("Failed to marshal message", err)
	}

	producerMessage := kafka.Message{
		Key:   []byte(createdUser.ID.String()),
		Value: messageBytes,
	}
	err = s.producer.WriteMessages(ctx, producerMessage)
	if err != nil {
		// Log error but don't fail the operation since alarm was created successfully
		// You might want to add proper logging here
		return nil, errors.NewInternalError("Failed to send message", err)
	}
	return createdUser, nil
}

func (s *UserService) GetUsers(ctx context.Context, page int, limit int) (*types.PaginateResponse[models.User], error) {
	users, err := s.userRepo.GetUsers(ctx, page, limit)
	if err != nil {
		return nil, errors.NewDatabaseError("get users", err)
	}
	total, err := s.userRepo.GetUsersCount(ctx)
	totalPages := int(total) / limit
	if total%int64(limit) != 0 {
		totalPages++
	}

	if err != nil {
		return nil, errors.NewDatabaseError("get users count", err)
	}
	paginateResponse := &types.PaginateResponse[models.User]{
		Pagination: types.Pagination{
			TotalPages: int(totalPages),
			Page:       page,
			Limit:      limit,
		},
		Data: users,
	}
	return paginateResponse, nil
}
func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, errors.NewDatabaseError("get user by id", err)
	}
	return user, nil
}
func (s *UserService) VerifyAccount(ctx context.Context, token string) error {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return errors.NewBadRequestError("Invalid token")
	}
	user, err := s.userRepo.GetUserByID(ctx, claims.UserID)

	if err != nil {
		return errors.NewDatabaseError("can not get user by id", err)
	}
	user.IsActive = true
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return errors.NewDatabaseError("can not update user", err)
	}
	return nil
}

// isDuplicateEmailError checks if the error is due to duplicate email constraint
func isDuplicateEmailError(err error) bool {
	errStr := err.Error()
	return contains(errStr, "duplicate key value violates unique constraint") &&
		contains(errStr, "email")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
