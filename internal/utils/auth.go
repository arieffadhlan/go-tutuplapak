package utils

import (
	"os"
	"strconv"
	"time"
	"tutuplapak-user/internal/entities"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func HashPasswordBeli(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func ValidToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

// Tambahkan fungsi ini ke dalam file services/auth.go

// generateJWTToken membuat token JWT berdasarkan data user
func GenerateJWTToken(user entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// Tambahkan email jika tersedia
	if user.Email != nil {
		claims["email"] = user.Email
	}

	// Tambahkan phone jika tersedia
	if user.Phone != nil {
		claims["phone"] = user.Phone
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateJWTTokenBeli(user entities.UserBeli) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"isAdmin":  user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
