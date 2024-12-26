package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	HashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(HashedPassword)
}