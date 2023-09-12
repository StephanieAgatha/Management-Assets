package model

// user details model
type UserCredentials struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UserLoginRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	IsActive bool   `json:"is_active"`
}

type UserLoginOTPRequest struct {
	OTP   int    `json:"otp"`
	Email string `json:"email" validate:"required,email"`
}
