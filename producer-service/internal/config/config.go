package config

import (
	"os"
)

type Config struct {
	KafkaBrokers string
	RedisAddr    string
	Port         string
}

func Load() Config {
	return Config{
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:9092"),
		RedisAddr:    getenv("REDIS_ADDR", "localhost:6379"),
		Port:         getenv("PORT", "8082"),
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
