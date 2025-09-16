package dto

type CreateUserDto struct {
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
	Role      string `json:"role" validate:"required,role"`
	PremiseID string `json:"premise_id" validate:"required,uuid"`
}
