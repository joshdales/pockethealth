package util

import (
	"fmt"
	"image/png"
	"os"
	"pockethealth/dicom/db"

	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

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

		imgFile, err := os.Create(fmt.Sprintf("images/%s.png", imgId))
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
