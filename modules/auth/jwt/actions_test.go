package jwt

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	checkGodotWorks(t)

	token, err := GenerateJWT("user1")
	if err != nil {
		log.Println(err)
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestParseJWTToken(t *testing.T) {
	checkGodotWorks(t)

	token, _ := GenerateJWT("user1")

	username, err := ParseJWTToken(token)

	assert.Nil(t, err)
	assert.Equal(t, username, "user1")
}

func checkGodotWorks(t *testing.T) {
	godotenv.Load("./.env")

	assert.Equal(t, os.Getenv("JWT_SECRET"), "test")
	assert.Equal(t, JWTConfig().SecretKey, "test")
	assert.Equal(t, os.Getenv("JWT_EXPIRATION"), "1")
}
