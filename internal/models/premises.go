package models

import "github.com/google/uuid"

type Premise struct {
	Base
	Name            string     `json:"name"`
	Address         string     `json:"address"`
	ParentPremiseID *uuid.UUID `json:"parent_premise_id,omitempty"`
	ParentPremise   *Premise   `json:"parent_premise,omitempty" gorm:"foreignKey:ParentPremiseID"`
}
