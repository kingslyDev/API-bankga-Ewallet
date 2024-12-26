package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Claims untuk JWT
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// LoadEnv memuat variabel lingkungan dari file .env
func LoadEnv() error {
	return godotenv.Load() // Memuat file .env
}

// GenerateJWT membuat token JWT
func GenerateJWT(email string) (string, time.Time, error) {
	if err := LoadEnv(); err != nil {
		return "", time.Time{}, fmt.Errorf("error loading .env file")
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		return "", time.Time{}, fmt.Errorf("JWT_SECRET_KEY is not set in .env")
	}

	expirationTime := time.Now().Add(24 * time.Hour) // Token berlaku selama 24 jam
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Membuat token dengan klaim yang telah ditentukan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan kunci JWT
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ParseJWT memverifikasi dan mem-parsing token JWT
func ParseJWT(tokenString string) (*Claims, error) {
	if err := LoadEnv(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtKey) == 0 {
		return nil, fmt.Errorf("JWT_SECRET_KEY is not set in .env")
	}

	// Mendeklarasikan klaim untuk menyimpan data yang diparse dari token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return claims, nil
}
