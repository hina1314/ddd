package validation_test

import (
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hina1314/ddd/internal/api/handler/dto"
	"github.com/hina1314/ddd/internal/api/validation"
	"github.com/stretchr/testify/require"
)

func TestPasswordValidatorAppliesToEveryUserRequest(t *testing.T) {
	v := validator.New()
	require.NoError(t, v.RegisterValidation("phone", validation.PhoneValidator))
	require.NoError(t, v.RegisterValidation("password", validation.PasswordValidator))

	overLimit := strings.Repeat("密", 25) // 25 characters, 75 UTF-8 bytes.
	tests := []struct {
		name    string
		request any
	}{
		{"registration", dto.CreateUserRequest{Type: 1, Password: overLimit}},
		{"login", dto.LoginUserRequest{Type: 1, Password: overLimit}},
		{"profile update", dto.UpdateUserRequest{Password: &overLimit}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.request)
			require.Error(t, err)
			validationErrors, ok := err.(validator.ValidationErrors)
			require.True(t, ok)
			require.Equal(t, "password", validationErrors[0].Tag())
		})
	}
}
