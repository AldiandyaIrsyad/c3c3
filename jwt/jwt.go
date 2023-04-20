package jwt

import (
	"fmt"
	"time"

	"github.com/aldiandyaIrsyad/c3c2/models"
	"github.com/dgrijalva/jwt-go"
)

const SecretKey = "your_secret_key"

func GenerateToken(username string, role models.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %v", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &models.User{
			Username: claims["username"].(string),
			Role:     models.Role(claims["role"].(string)),
		}
		return user, nil
	}

	return nil, err
}
