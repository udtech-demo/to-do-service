package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPwd(s string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashed)
}

func ComparePwd(hashed string, normal string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(normal))
}

func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return errors.New("invalid validation error")
		}

		for _, err := range err.(validator.ValidationErrors) {

			switch err.ActualTag() {
			case "required":
				msg := fmt.Sprintf("Required parameters not passed (%s)", err.Field())
				return errors.New(msg)

			default:
				msg := fmt.Sprintf("Parameters incorrectly formatted or out of range (%s)", err.Field())
				return errors.New(msg)
			}
		}
	}

	return nil
}
