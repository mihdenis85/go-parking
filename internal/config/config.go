package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LoggingFormat        string
	LoggingDateFormat    string
	LoggingLevel         string
	MongoDBURL           string
	AppTitle             string
	DBName               string
	ParkingServiceAPIKey string
	ParkingSlotsCount    int
	ServerPort           string
}

var Settings *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	Settings = &Config{
		LoggingFormat:        getEnv("LOGGING_FORMAT", "%(asctime)s,%(msecs)d %(levelname)-8s [%(pathname)s:%(lineno)d in function %(funcName)s] %(message)s"),
		LoggingDateFormat:    getEnv("LOGGING_DATE_FORMAT", "2006-01-02 15:04:05"),
		LoggingLevel:         getEnv("LOGGING_LEVEL", "INFO"),
		MongoDBURL:           getEnv("MONGODB_URL", "mongodb://mongodb:27017"),
		AppTitle:             getEnv("APP_TITLE", "Parking Service"),
		DBName:               getEnv("DB_NAME", "ParkingService"),
		ParkingServiceAPIKey: getEnvRequired("PARKING_SERVICE_API_KEY"),
		ParkingSlotsCount:    getEnvAsInt("PARKING_SLOTS_COUNT", 52),
		ServerPort:           getEnv("SERVER_PORT", "8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s is not set", key)
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
