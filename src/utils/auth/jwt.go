package auth

import (
	"errors"
	"os"
	"pocok/src/utils/models"

	"github.com/golang-jwt/jwt"
)

func CreateToken(orgId string) (string, error) {
	jwtKey := []byte(os.Getenv("jwtKey"))
	expiry := jwt.TimeFunc().Unix() + 86400*2 // 2 days

	claims := models.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
		},
		JWTCustomClaims: models.JWTCustomClaims{
			OrgId: orgId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, createTokenError := token.SignedString(jwtKey)
	return tokenString, createTokenError
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ParseToken(tokenString string) (*models.JWTClaims, error) {
	jwtKey := []byte(os.Getenv("jwtKey"))

	token, jwtParseError := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if jwtParseError != nil {
		return nil, jwtParseError
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return claims, errors.New("invalid JWT Token")
}
