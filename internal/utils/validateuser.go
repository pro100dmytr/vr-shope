package utils

import (
	"errors"
	"vr-shope/internal/models/services"
)

func ValidateUser(user *services.User) error {
	if user.Login == "" {
		return errors.New("login is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
