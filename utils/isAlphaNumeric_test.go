package utils_test

import (
	"testing"

	"github.com/dockrelix/dockrelix-backend/utils"
)

func TestIsAlphaNumeric(t *testing.T) {
	t.Run("Valid alphanumeric", func(t *testing.T) {
		input := "dockrelix123"
		if !utils.IsAlphaNumeric(input) {
			t.Errorf("Alphanumeric validation failed for %s", input)
		}
	})

	symbols := []string{"@", " ", "-", "_", "!", "#", "$", "%", "^", "&", "*", "(", ")", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", ".", "?", "/", "`", "~"}
	for _, symbol := range symbols {
		t.Run("Invalid alphanumeric with symbol "+symbol, func(t *testing.T) {
			input := "dockrelix" + symbol + "123"
			if utils.IsAlphaNumeric(input) {
				t.Errorf("Alphanumeric validation failed for %s", input)
			}
		})
	}
}
