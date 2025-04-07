package dto

import "time"

type CreateUserRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Username string `json:"username" validate:"required,alphanumunicode"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"email"`
}

type UserResponse struct {
	Phone       string    `json:"phone"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
