package config

import (
	"github.com/joho/godotenv"
	"time"
)

func Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig interface {
	Address() string
	AccessClientAddress() string
	RateLimit() (int, time.Duration)
}

type PGConfig interface {
	DSN() string
}
