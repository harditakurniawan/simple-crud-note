package utils

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func CheckContainSpecialChar(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	re := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
	return re.MatchString(password)
}

func toJSONString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
