package assemble

import (
	"study/internal/api/handler/dto"
	"study/token"
)

type UpdateUserCommand struct {
	ID       int64
	Phone    string
	Email    string
	Username string
	Password string
}

func NewUpdateUserCommand(user dto.UpdateUserRequest, payload *token.Payload) *UpdateUserCommand {
	return &UpdateUserCommand{
		ID:       payload.UserId,
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
}
