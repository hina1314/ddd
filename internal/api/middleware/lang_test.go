package middleware

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectLocale(t *testing.T) {
	supported := map[string]string{"en": "en", "zh": "zh"}

	tests := []struct {
		name   string
		header string
		want   string
	}{
		{name: "exact", header: "zh", want: "zh"},
		{name: "regional", header: "zh-CN, en;q=0.8", want: "zh"},
		{name: "fallback", header: "fr-FR", want: "en"},
		{name: "quality suffix", header: "fr;q=0.9, en;q=0.8", want: "en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, selectLocale(tt.header, "en", supported))
		})
	}
}
