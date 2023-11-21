package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Host         string
	DatabaseURI  string
	LoggingLevel string
	AccrualHost  string
	SecretKey    string
	Multiplier   int
}

var AppConfig Config

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

func New() *Config {
	loadFlagConfig(&AppConfig)
	loadEnvConfig(&AppConfig)

	wd, _ := os.Getwd()
	envFile, err := godotenv.Read(filepath.Join(wd, "config.env"))
	if err != nil {
		log.Fatal("Error loading config.env file")
	}

	AppConfig.SecretKey = envFile["SECRET_KEY"]
	AppConfig.Multiplier, _ = strconv.Atoi(envFile["MULTIPLIER"])

	return &AppConfig
}

func GetConfig() *Config {
	return &AppConfig
}
