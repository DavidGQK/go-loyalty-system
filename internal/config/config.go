package config

import (
	"flag"
	"os"
)

type Config struct {
	Host         string
	DatabaseURI  string
	LoggingLevel string
	AccrualHost  string
}

func loadFlagConfig(AppConfig *Config) {
	flag.StringVar(&AppConfig.Host, "a", "localhost:8080", "url where server runs on")
	flag.StringVar(&AppConfig.DatabaseURI, "d", "", "data for db connection")
	flag.StringVar(&AppConfig.LoggingLevel, "l", "info", "logging level")
	flag.StringVar(&AppConfig.AccrualHost, "r", "localhost:8080", "accrual host")

	flag.Parse()
}

func loadEnvConfig(AppConfig *Config) {
	if envHost := os.Getenv("RUN_ADDRESS"); envHost != "" {
		AppConfig.Host = envHost
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		AppConfig.DatabaseURI = envDatabaseURI
	}

	if envLoggingLevel := os.Getenv("LOG_LEVEL"); envLoggingLevel != "" {
		AppConfig.LoggingLevel = envLoggingLevel
	}

	if envAccrualHost := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualHost != "" {
		AppConfig.AccrualHost = envAccrualHost
	}
}

func GetConfig() *Config {
	var AppConfig Config

	loadFlagConfig(&AppConfig)
	loadEnvConfig(&AppConfig)

	return &AppConfig
}
