package models

import "github.com/google/uuid"

type UserPremise struct {
	Base
	UserID    uuid.UUID `json:"user_id"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PremiseID uuid.UUID `json:"premise_id"`
	Premise   *Premise  `json:"premise,omitempty" gorm:"foreignKey:PremiseID"`
}
