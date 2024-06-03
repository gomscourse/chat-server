package config

import (
	"github.com/joho/godotenv"
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
}

type PGConfig interface {
	DSN() string
}
