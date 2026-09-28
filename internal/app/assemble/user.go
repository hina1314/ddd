package assemble

import (
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/token"
)

type UpdateUserCommand struct {
	ID       int64
	Phone    *string
	Email    *string
	Username *string
	Password *string
}

func NewUpdateUserCommand(user dto.UpdateUserRequest, payload *token.Payload) *UpdateUserCommand {
	return &UpdateUserCommand{
		ID:       payload.UserID,
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
}
