package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Loadenv() {
	var err error
	zap.S().Info("loading .env...")

	if os.Getenv("APP_ENV") == "test" {
		err = godotenv.Load("../../../.env")
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		zap.S().Warn("Error loading .env", err.Error())
	}
	checkenv()
}

func checkenv() {
	if os.Getenv("APP_ENV") == "" {
		zap.S().Fatal("APP_ENV required")
	}
	if os.Getenv("MONGO_URI") == "" {
		zap.S().Fatal("MONGO_URI required")
	}
	if os.Getenv("GCP_PROJECTID") == "" {
		zap.S().Fatal("GCP_PROJECTID required")
	}
	if os.Getenv("GCP_BUCKETNAME") == "" {
		zap.S().Fatal("GCP_BUCKETNAME required")
	}
	if os.Getenv("GCP_PREFIX") == "" {
		zap.S().Fatal("GCP_PREFIX required")
	}
	if os.Getenv("REDIS_ENDPOINT") == "" {
		zap.S().Fatal("REDIS_ENDPOINT required")
	}
	if os.Getenv("REDIS_PORT") == "" {
		zap.S().Fatal("REDIS_PORT required")
	}
	if os.Getenv("REDIS_DB_PASSWORD") == "" && os.Getenv("APP_ENV") != "development" && os.Getenv("APP_ENV") != "test" {
		zap.S().Fatal("REDIS_DB_PASSWORD required")
	}
	if os.Getenv("JWT_SECRET") == "" {
		zap.S().Fatal("JWT_SECRET required")
	}
	if os.Getenv("ALLOW_ORIGINS_ENDPOINT") == "" {
		zap.S().Fatal("ALLOW_ORIGINS_ENDPOINT required")
	}

}
