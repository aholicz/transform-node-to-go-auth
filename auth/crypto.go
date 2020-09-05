package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), viper.GetInt("passwordLeght"))
	return string(bytes), err
}

func validatePassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func generateJwtToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("jwt.timeout")).Unix(),
		},
	})

	signed, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return signed, err
}
