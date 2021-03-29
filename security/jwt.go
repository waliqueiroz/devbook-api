package security

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/waliqueiroz/devbook-api/config"
)

// CreateToken generates a new json web token for a given user
func CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userID"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	return token.SignedString(config.SecretKey)
}
