package auth

import (
	"backend/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CheckLoginValidity(userData *models.User, loginData *models.LoginUserData) (int64, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(userData.PasswordHash),
		[]byte(loginData.Password),
	)

	if err != nil {
		return 0, err
	}

	return userData.ID, nil
}

func GenerateJWT(userID int64) (string, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if len(jwtSecret) == 0 {
		log.Println("Empty secret")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)

	return tokenStr, err
}

func CheckTokenValidity(tokenStr string) (*jwt.Token, error) {
	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	if len(jwtSecret) == 0 {
		log.Println("Empty Secret")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, match := token.Method.(*jwt.SigningMethodHMAC); !match {
			return nil, jwt.ErrTokenMalformed
		}

		return jwtSecret, nil
	})

	return token, err
}
