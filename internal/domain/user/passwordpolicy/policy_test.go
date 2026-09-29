package passwordpolicy

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValid(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"seven ASCII characters", "1234567", false},
		{"eight ASCII characters", "12345678", true},
		{"seven multibyte characters", "密码密码密码密", false},
		{"eight multibyte characters", "密码密码密码密码", true},
		{"exactly 72 bytes", strings.Repeat("a", MaxBytes), true},
		{"more than 72 ASCII bytes", strings.Repeat("a", MaxBytes+1), false},
		{"more than 72 multibyte bytes", strings.Repeat("密", 25), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.valid, Valid(tt.password))
		})
	}
}
