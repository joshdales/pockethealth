package util

import (
	"fmt"
	"image/png"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"

	"pockethealth/dicom/db"
)

// In an actual production env I'm sure this would upload to your Object Storage, but I'm just storing locally
func UploadImage(file multipart.File, patientId string) (string, error) {
	var fileData []byte
	file.Read(fileData)

	imageId := uuid.New()
	filePath := fmt.Sprintf("%s/%s.dcm", db.StorageLocation, imageId)

	err := os.WriteFile(filePath, fileData, 0660)
	if err != nil {
		return "", fmt.Errorf("Could not save file: %e", err)
	}

	return imageId.String(), nil
}

// Convert the DICOM image to png and upload it (or just store it locally in this case)
func ConvertDicomToPngAndUpload(dataset *dicom.Dataset, dicomImgId string, patientId string) error {
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return fmt.Errorf("Could not get pixel data for image: %e", err)
	}

	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)

	for _, frame := range pixelDataInfo.Frames {
		img, err := frame.GetImage()
		if err != nil {
			return fmt.Errorf("Could not get image frame: %e", err)
		}

		imgId := uuid.New()

		imgFile, err := os.Create(fmt.Sprintf("%s/%s.png", db.StorageLocation, imgId))
		if err != nil {
			return fmt.Errorf("Could not create png image file: %e", err)
		}

		err = png.Encode(imgFile, img)
		if err != nil {
			return fmt.Errorf("Could not encode png image: %e", err)
		}

		err = imgFile.Close()
		if err != nil {
			return err
		}

		db.CreatePngImage(imgId.String(), dicomImgId, patientId)
	}

	return nil
}
