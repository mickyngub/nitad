package gcp

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/utils"
	"go.uber.org/zap"
)

type Uploader interface {
	UploadFiles(ctx context.Context, files []*multipart.FileHeader, collectionName string) ([]string, errors.CustomError)
	UploadFile(ctx context.Context, file *multipart.FileHeader, collectionName string) (string, errors.CustomError)
	DeleteFiles(ctx context.Context, URLs []string)
	DeleteFile(ctx context.Context, URL string)
}

type uploader struct {
	cl         *storage.Client
	bucketName string
}

func Init() uploader {
	if os.Getenv("APP_ENV") != "production" {
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "google-credentials.json")
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		zap.S().Fatal("Failed to create gcp client: ", err.Error())
	}

	return uploader{
		cl:         client,
		bucketName: os.Getenv("GCP_BUCKETNAME"),
	}

}

// UploadFiles uploads multiple files
// just loop through all files and pass to
// the UploadFile function one by one
// return an array of filenames used to concat to the gcp baseURLs
func (u uploader) UploadFiles(ctx context.Context, files []*multipart.FileHeader, collectionName string) ([]string, errors.CustomError) {
	filenames := []string{}
	for _, file := range files {
		//TODO: channel this
		filename, err := u.UploadFile(ctx, file, collectionName)
		if err != nil {
			return filenames, errors.NewBadRequestError(err.Error())
		}
		filenames = append(filenames, filename)
	}
	return filenames, nil
}

// UploadFile uploads a single file
// return a filename used to concat to the gcp baseURLs
// ex: reports/28-Feb-2022-18:57:15-dummyReport.pdf
// 		 images/28-Feb-2022-18:11:11-dummyImage.png
func (u uploader) UploadFile(ctx context.Context, file *multipart.FileHeader, collectionName string) (string, errors.CustomError) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var filename = ""
	blobFile, err := file.Open()
	if err != nil {
		return filename, errors.NewBadRequestError(err.Error())
	}
	defer blobFile.Close()

	filename, filetype := utils.GetUniqueFilename(file.Filename)

	// Upload an object with storage.Writer to the uploadPath
	uploadPath := collectionName + "/" + filetype + "/" + filename
	wc := u.cl.Bucket(u.bucketName).Object(uploadPath).NewWriter(ctx)
	if _, err := io.Copy(wc, blobFile); err != nil {
		// zap.S().Warn("GCP io.Copy: ", err.Error())
		return filename, errors.NewInternalServerError("GCP io.Copy: " + err.Error())
	}
	if err := wc.Close(); err != nil {
		// zap.S().Warn("GCP Writer.Close: ", err.Error())
		return filename, errors.NewInternalServerError("GCP Writer.Close: " + err.Error())
	}

	return uploadPath, nil
}

// DeleteFiles delete multiple files by looping through each one
// and pass through the DeleteFile function
func (u uploader) DeleteFiles(ctx context.Context, filepaths []string) {
	for _, filepath := range filepaths {
		//TODO: channel this
		u.DeleteFile(ctx, filepath)
	}
}

// DeleteFile removes a specified file.
func (u uploader) DeleteFile(ctx context.Context, filepath string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := u.cl.Bucket(u.bucketName).Object(filepath)
	if err := o.Delete(ctx); err != nil {
		zap.S().Warn("gcp deletion error, file= ", filepath, " ", err.Error())
	}
	zap.S().Info("Deleting..", filepath)
}

func GetURL(filepath string) string {
	return os.Getenv("GCP_PREFIX") + "/" + os.Getenv("GCP_BUCKETNAME") + "/" + filepath
}

func GetFilepath(URL string) string {
	arr := strings.Split(URL, "/")
	filename := strings.Join(arr[len(arr)-3:], "/")
	return filename
}
