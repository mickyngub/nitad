package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
)

type ClientUploader struct {
	cl         *storage.Client
	bucketName string
	apiPrefix  string
}

var uploader *ClientUploader

func Init() {
	if os.Getenv("APP_ENV") != "production" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "google-credentials.json")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create gcp client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: os.Getenv("GCP_BUCKETNAME"),
		apiPrefix:  os.Getenv("GCP_API_PREFIX"),
	}
}

// UploadFile uploads an object
func UploadFile(f multipart.File, object string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := uploader.cl.Bucket(uploader.bucketName).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return fmt.Errorf("GCP io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("GCP Writer.Close: %v", err)
	}

	return nil
}

func UploadImages(files []*multipart.FileHeader, collectionName string) ([]string, errors.CustomError) {
	urls := []string{}
	for _, file := range files {
		blobFile, err := file.Open()
		if err != nil {
			return urls, errors.NewBadRequestError(err.Error())
		}
		defer blobFile.Close()
		filename := functions.GetUniqueFilename(file.Filename)

		//TODO: channel this
		err = UploadFile(blobFile, collectionName+"/"+filename)
		if err != nil {
			return urls, errors.NewBadRequestError(err.Error())
		}
		urls = append(urls, filename)
	}

	return urls, nil
}

func DeleteImages(imageURLS []string, collectionName string) errors.CustomError {
	for _, url := range imageURLS {
		filepath := collectionName + "/" + url

		//TODO: channel this
		err := DeleteFile(filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteFile removes specified object.
func DeleteFile(object string) errors.CustomError {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	o := uploader.cl.Bucket(uploader.bucketName).Object(object)
	if err := o.Delete(ctx); err != nil {
		return errors.NewInternalServerError("gcp deletion error, " + err.Error())
	}

	return nil
}
