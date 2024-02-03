package memphis

import (
	"os"

	"github.com/memphisdev/memphis.go"
)

type Config struct {
	Hostname        string
	ApplicationUser string
	Password        memphis.Option
}

func MemphisConfig() *Config {
	return &Config{
		Hostname:        os.Getenv("MEMPHIS_HOSTNAME"),
		ApplicationUser: os.Getenv("MEMPHIS_APP_USER"),
		Password:        memphis.Password(os.Getenv("MEMPHIS_APP_PASSWORD")),
	}
}
