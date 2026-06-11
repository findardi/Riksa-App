package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=6,max=255"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Password string `json:"password" validate:"required"`
}
