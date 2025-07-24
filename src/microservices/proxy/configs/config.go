package config

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type RouteConfig struct {
	PathPrefix string `yaml:"path_prefix"`
	Target     string `yaml:"target"`
}

type Config struct {
	Routes           []RouteConfig `yaml:"routes"`
	MigrationPercent int
	GradualMigration bool
}

func Load() *Config {

	configFile, err := os.Open("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to open configs/config.yaml: %v", err)
		panic(err)
	}
	defer configFile.Close()

	var config Config

	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config.yaml: %v", err)
		panic(err)
	}
	migrationPerc := GetEnv("MOVIES_MIGRATION_PERCENT", "100")
	gradualMigration := GetEnv("GRADUAL_MIGRATION", "true")

	if gradualMigration == "true" {
		config.GradualMigration = true
	} else {
		config.GradualMigration = false
	}
	config.MigrationPercent, _ = strconv.Atoi(migrationPerc)
	log.Printf("Config loadid: %v", config)
	return &config
}

func GetEnv(key, defaulValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaulValue
}
