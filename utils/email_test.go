package utils_test

import (
	"testing"

	"github.com/dockrelix/dockrelix-backend/utils"
)

func TestEmailValidation(t *testing.T) {
	t.Run("Valid email", func(t *testing.T) {
		email := "contact@dockrelix.org"
		if !utils.IsEmailValid(email) {
			t.Errorf("Email validation failed for %s", email)
		}
	})

	t.Run("Invalid email", func(t *testing.T) {
		email := "contact@dockrelix"
		if utils.IsEmailValid(email) {
			t.Errorf("Email validation failed for %s", email)
		}
	})
}
