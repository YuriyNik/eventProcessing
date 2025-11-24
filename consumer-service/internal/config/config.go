package config

import "os"

type Config struct {
	KafkaBrokers string
	Topic        string
	GroupID      string
}

func Load() Config {
	return Config{
		KafkaBrokers: getenv("KAFKA_BROKERS", "localhost:19092"),
		Topic:        getenv("TOPIC", "events"),
		GroupID:      getenv("GROUP_ID", "consumer-group-1"),
	}
}

func getenv(k, d string) string {
	v := os.Getenv(k)
	if v == "" {
		return d
	}
	return v
}
