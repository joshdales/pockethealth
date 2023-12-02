package model

import (
	"time"
)

type PngImage struct {
	Id           string    `json:"id"`
	DicomImageId string    `json:"dicom_image_id"`
	PatientId    string    `json:"patient_id"`
	StorageUrl   string    `json:"storage_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
