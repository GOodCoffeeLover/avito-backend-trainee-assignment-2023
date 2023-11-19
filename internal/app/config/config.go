package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	appPortEnvKey               = "APP_PORT"
	postgresHostEnvKey          = "POSTGRES_HOST"
	postgresPortEnvKey          = "POSTGRES_PORT"
	postgresUserEnvKey          = "POSTGRES_USER"
	postgresPassEnvKey          = "POSTGRES_PASS"
	postgresExtraConnOptsEnvKey = "POSTGRES_EXTRA_CONNECT_OPTIONS"
)

type Config struct {
	Port       uint
	ConnString string
}

// TODO: make config from yaml
func New() (*Config, error) {
	port, err := strconv.Atoi(GetEnvWithDefault(appPortEnvKey, "7001"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse port: %w", err)
	}
	return &Config{
		Port: uint(port),
		ConnString: fmt.Sprintf("host=%v port=%v user=%v password=%v",
			GetEnvWithDefault(postgresHostEnvKey, "localhost"), GetEnvWithDefault(postgresPortEnvKey, "5432"),
			GetEnvWithDefault(postgresUserEnvKey, "postgres"), GetEnvWithDefault(postgresPassEnvKey, "postgres")),
	}, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
