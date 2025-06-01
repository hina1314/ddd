package dto

import "time"

type CreateUserRequest struct {
	Type     int8   `json:"type" validate:"required,oneof=1 2"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Phone    string `json:"phone" validate:"omitempty,phone"`
	Email    string `json:"email" validate:"omitempty,email"`
	Username string `json:"username" validate:"omitempty,min=1,max=32"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

type UserResponse struct {
	Phone       string    `json:"phone"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Type     int8   `json:"type" validate:"required,oneof=1 2"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=6"`
}
