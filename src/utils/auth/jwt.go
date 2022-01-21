package auth

import (
	"errors"
	"os"
	"pocok/src/utils/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJwt(orgId string) (string, error) {
	jwtKey := os.Getenv("jwtKey")
	claims := models.JWTClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400*2,
		},
		models.JWTPayload{
			OrgId: orgId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ParseJwt(tokenString string) (*models.JWTClaims, error) {
	jwtKey := os.Getenv("jwtKey")

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return claims, errors.New("invalid JWT Token")
}
