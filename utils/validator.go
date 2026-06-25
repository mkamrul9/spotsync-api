package utils

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator package for Echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate method matches Echo's Validator interface
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
