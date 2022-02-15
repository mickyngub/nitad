package gcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/functions"
)

var (
	projectID  = os.Getenv("GCP_PROJECTID")
	bucketName = os.Getenv("GCP_BUCKETNAME")
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func Init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "google-credentials.json") // FILL IN WITH YOUR FILE PATH
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
	}

}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, object string, uploadPath string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
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
		filename := functions.GetUniqueFilename(file.Filename)

		//TODO: channel this
		err = uploader.UploadFile(blobFile, filename, fmt.Sprintf("%s/", collectionName))
		if err != nil {
			return urls, errors.NewBadRequestError(err.Error())
		}

		urls = append(urls, fmt.Sprintf("https://storage.cloud.google.com/nitad/%s/%s", collectionName, filename))
	}
	return urls, nil
}

func DeleteImages(imageURLS []string, collectionName string) errors.CustomError {
	for _, url := range imageURLS {
		urlSlice := strings.Split(url, "/")
		filepath := fmt.Sprintf("%s/%s", collectionName, urlSlice[len(urlSlice)-1])

		//TODO: channel this
		err := DeleteFile(filepath)
		if err != nil {
			return errors.NewBadRequestError(err.Error())
		}
	}
	return nil
}

// deleteFile removes specified object.
func DeleteFile(object string) errors.CustomError {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucketName).Object(object)
	if err := o.Delete(ctx); err != nil {
		return errors.NewBadRequestError(err.Error())
	}

	return nil
}
