package util

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
)

// In an actual production env I'm sure this would upload to your Object Storage, but I'm just storing locally
func UploadImage(file multipart.File, patientId string) (string, error) {
	var fileData []byte
	file.Read(fileData)

	imageId := uuid.New()
	filePath := fmt.Sprintf("images/%s.dcm", imageId)

	err := os.Mkdir("images", 0750)
	if err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("Couldn't create images folder: %e", err)
	}

	err = os.WriteFile(filePath, fileData, 0660)
	if err != nil {
		return "", fmt.Errorf("Could not save file: %e", err)
	}

	return imageId.String(), nil
}
