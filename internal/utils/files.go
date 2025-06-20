package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/rs/xid"
)

func SaveFile(file *multipart.FileHeader) (string, error) {

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Destination
	var fileName string = xid.New().String() + filepath.Ext(file.Filename)
	dst, err := os.Create(os.Getenv("PATH_TO_UPLOAD") + fileName)
	if err != nil {
		return "", err
	}

	// Copy the uploaded content to the destination file.
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return fileName, nil
}
