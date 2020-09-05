package auth

import (
	"fmt"
)

const (
	RequiredEmailMessage    = "email is required"
	RequiredPasswordMessage = "password is required"
)

func isMissingRequired(authProfile *AuthProfile) (bool, string) {
	var (
		msg     string
		missing bool
	)

	if authProfile.Email == "" {
		msg += RequiredEmailMessage
		missing = true
	}

	if authProfile.Password == "" {
		msg += fmt.Sprintf(", %s", RequiredPasswordMessage)
		missing = true
	}

	return missing, msg
}
