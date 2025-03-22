package utils_test

import (
	"testing"

	"github.com/dockrelix/dockrelix-backend/utils"
)

func TestSanitizeInput(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		input := "dockrelix"
		output := utils.SanitizeInput(input)
		if input != output {
			t.Errorf("Sanitization failed for %s", input)
		}
	})

	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"Single quote", "dock'relix", "dockrelix"},
		{"Double quote", "dock\"relix", "dockrelix"},
		{"Semicolon", "dock;relix", "dockrelix"},
		{"Double dash", "dock--relix", "dockrelix"},
		{"Block comment start", "dock/*relix", "dockrelix"},
		{"Block comment end", "dock*/relix", "dockrelix"},
		{"XP_", "dockxp_relix", "dockrelix"},
		{"Slash", "dock/relix", "dockrelix"},
		{"Multiple", "dock'rel;ix--", "dockrelix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := utils.SanitizeInput(tt.input)
			if output != tt.output {
				t.Errorf("Sanitization failed for %s: expected %s, got %s", tt.input, tt.output, output)
			}
		})
	}
}
