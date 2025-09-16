package repositories

import (
	"context"
	"scs-user/internal/models"

	"gorm.io/gorm"
)

type UserPremiseRepository struct {
	db *gorm.DB
}

func NewUserPremiseRepository(db *gorm.DB) *UserPremiseRepository {
	return &UserPremiseRepository{db: db}
}

func (r *UserPremiseRepository) AssignPremises(ctx context.Context, userPremise *models.UserPremise) (error) {
	if err := r.db.WithContext(ctx).Create(userPremise).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserPremiseRepository) CheckExist(ctx context.Context, userPremise *models.UserPremise) (bool, error) {
	existingUserPremise := &models.UserPremise{}
	if err := r.db.WithContext(ctx).Where("guard_id = ? AND premise_id = ?", userPremise.UserID, userPremise.PremiseID).First(existingUserPremise).Error; err == nil {
		return true, nil
	}
	return false, nil
}
