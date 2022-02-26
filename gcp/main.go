package gcp

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"go.uber.org/zap"
)

type ClientUploader interface {
	UploadFile(ctx context.Context, f multipart.File, object string) error
	UploadImages(ctx context.Context, files []*multipart.FileHeader, collectionName string) ([]string, errors.CustomError)
	DeleteFile(ctx context.Context, object string) errors.CustomError
	DeleteImages(ctx context.Context, imageURLS []string, collectionName string) errors.CustomError
}

type clientUploader struct {
	cl         *storage.Client
	bucketName string
	apiPrefix  string
}

func Init() ClientUploader {
	if os.Getenv("APP_ENV") != "production" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "google-credentials.json")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		zap.S().Fatal("Failed to create gcp client: ", err.Error())
	}

	return &clientUploader{
		cl:         client,
		bucketName: os.Getenv("GCP_BUCKETNAME"),
		apiPrefix:  os.Getenv("GCP_API_PREFIX"),
	}

}

// UploadFile uploads an object
func (uploader *clientUploader) UploadFile(ctx context.Context, f multipart.File, object string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := uploader.cl.Bucket(uploader.bucketName).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		zap.S().Fatal("GCP io.Copy: ", err.Error())
	}
	if err := wc.Close(); err != nil {
		zap.S().Fatal("GCP Writer.Close: ", err.Error())
	}

	return nil
}

func (uploader *clientUploader) UploadImages(ctx context.Context, files []*multipart.FileHeader, collectionName string) ([]string, errors.CustomError) {
	urls := []string{}
	for _, file := range files {
		blobFile, err := file.Open()
		if err != nil {
			return urls, errors.NewBadRequestError(err.Error())
		}
		defer blobFile.Close()
		filename := utils.GetUniqueFilename(file.Filename)

		//TODO: channel this
		err = uploader.UploadFile(ctx, blobFile, collectionName+"/"+filename)
		if err != nil {
			return urls, errors.NewBadRequestError(err.Error())
		}
		urls = append(urls, filename)
	}

	return urls, nil
}

func (uploader *clientUploader) DeleteImages(ctx context.Context, imageURLS []string, collectionName string) errors.CustomError {
	for _, url := range imageURLS {
		filepath := collectionName + "/" + url

		//TODO: channel this
		err := uploader.DeleteFile(ctx, filepath)
		if err != nil {
			return err
		}
	}
	return nil
}

// deleteFile removes specified object.
func (uploader *clientUploader) DeleteFile(ctx context.Context, object string) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := uploader.cl.Bucket(uploader.bucketName).Object(object)
	if err := o.Delete(ctx); err != nil {
		return errors.NewInternalServerError("gcp deletion error, " + err.Error())
	}

	return nil
}
