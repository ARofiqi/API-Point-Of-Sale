package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"ConfirmPassword" validate:"required,min=6"`
}

type RegisterResponse struct {
	ID    interface{} `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
