package config

import (
	"fmt"
	"os"
)

const (
	PostgresHost     = "POSTGRES_HOST"
	PostgresUsername = "POSTGRES_USERNAME"
	PostgresPassword = "POSTGRES_PASSWORD"
	PostgresDB       = "POSTGRES_DB"
)

var conf = map[string]string{}

func InitConfig() error {
	for _, field := range []string{PostgresHost, PostgresUsername, PostgresPassword, PostgresDB} {
		val, exists := os.LookupEnv(field)
		if !exists {
			return fmt.Errorf("field %s not found", field)
		}

		conf[field] = val
	}

	return nil
}

func Get(field string) (string, error) {
	val, ok := conf[field]
	if !ok {
		return "", fmt.Errorf("field %s not found", field)
	}

	return val, nil
}
