package utils

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func HashPassword(password string) string {
	HashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(HashedPassword)
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil 
}



func RandomString(n int) string {

    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	rand.Seed(uint64(time.Now().UnixNano()))

    b := make([]rune, n)

    for i := range b {

        b[i] = letters[rand.Intn(len(letters))]

    }

    return string(b)

}
