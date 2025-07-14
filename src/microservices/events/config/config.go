package config

import "os"

type Config struct {
	KafkaBroker       string
	KafkaUserTopic    string
	KafkaMovieTopic   string
	KafkaPaymentTopic string
	ServerPort        string
}

func Load() *Config {
	return &Config{
		KafkaBroker:       GetEnv("KAFKA_BROKERS", "localhost:9092"),
		ServerPort:        GetEnv("PORT", "8082"),
		KafkaUserTopic:    GetEnv("KAFKA_TOPIC_USERS", "user-events"),
		KafkaMovieTopic:   GetEnv("KAFKA_TOPIC_MOVIES", "movie-events"),
		KafkaPaymentTopic: GetEnv("KAFKA_TOPIC_PAYMENTS", "payment-events"),
	}
}

func GetEnv(key, defaulValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaulValue
}
