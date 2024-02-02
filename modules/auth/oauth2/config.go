package oauth2

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
	// LinkedInConfig
	// FacebookConfig
	// TwitterConfig
}

func GoogleConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("cannot load .env")
	}

	return &Config{
		GoogleLoginConfig: oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_CALLBACK_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{""}, // a mechanism to limit application access to user data
			Endpoint:     google.Endpoint,
		},
	}
}
