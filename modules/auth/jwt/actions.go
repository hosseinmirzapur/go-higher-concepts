package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	// set expiration time
	expirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	expTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// form the claim struct
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// create token from claim and sign it with secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTConfig().SecretKey))
}

func ParseJWTToken(token string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(jwtTkn *jwt.Token) (interface{}, error) {
		return []byte(JWTConfig().SecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid signature")
		}

		if err == jwt.ErrTokenExpired {
			return "", errors.New("token expired")
		}

		return "", err
	}

	if !tkn.Valid {
		return "", errors.New("unauthorized")
	}

	return claims.Username, nil
}
