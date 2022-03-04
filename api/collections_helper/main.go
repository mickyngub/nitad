package collections_helper

import (
	"context"
	"mime/multipart"

	"github.com/birdglove2/nitad-backend/errors"
	"github.com/birdglove2/nitad-backend/gcp"
)

func HandleUpdateSingleFile(u gcp.Uploader, c context.Context, newFile *multipart.FileHeader, oldFilename string, collectionName string) (string, errors.CustomError) {
	// delete old file
	u.DeleteFile(c, oldFilename, collectionName)

	// upload new file
	newUploadFilename, err := u.UploadFile(c, newFile, collectionName)
	if err != nil {
		// if uploading new file error, it might already be uploaded, so try deleting it as well
		u.DeleteFile(c, oldFilename, collectionName)
		return "", err
	}

	return newUploadFilename, nil
}
