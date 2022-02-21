package config

import (
	"fmt"
	"os"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/joho/godotenv"
)

func Loadenv() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error loading .env file!")
		}
	}
}

func Checkenv() errors.CustomError {
	if os.Getenv("APP_ENV") == "" {
		return errors.NewInternalServerError("APP_ENV required")
	}
	if os.Getenv("MONGO_URI") == "" {
		return errors.NewInternalServerError("MONGO_URI required")
	}
	if os.Getenv("GCP_PROJECTID") == "" {
		return errors.NewInternalServerError("GCP_PROJECTID required")
	}
	if os.Getenv("GCP_BUCKETNAME") == "" {
		return errors.NewInternalServerError("GCP_BUCKETNAME required")
	}
	// if os.Getenv("GCP_SERVICE_ACCOUNT") == "" {
	// 	return errors.NewInternalServerError("GCP_SERVICE_ACCOUNT required")
	// }
	if os.Getenv("GCP_API_PREFIX") == "" {
		return errors.NewInternalServerError("GCP_API_PREFIX required")
	}
	if os.Getenv("GCP_WORKLOAD_IDENTITY_PROVIDER") == "" {
		return errors.NewInternalServerError("GCP_WORKLOAD_IDENTITY_PROVIDER required")
	}
	if os.Getenv("REDIS_ENDPOINT") == "" {
		return errors.NewInternalServerError("REDIS_ENDPOINT required")
	}
	// if os.Getenv("REDIS_DB_PASSWORD") == "" {
	// 	return errors.NewInternalServerError("REDIS_DB_PASSWORD required")
	// }
	if os.Getenv("JWT_SECRET") == "" {
		return errors.NewInternalServerError("JWT_SECRET required")
	}

	return nil
}
