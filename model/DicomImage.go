package model

import (
	"time"

	"github.com/suyashkumar/dicom"
)

type DicomImage struct {
	Id               string        `json:"id"`
	UploadedByUserId string        `json:"uploaded_by_user_id"`
	PatientId        string        `json:"patient_id"`
	StorageUrl       string        `json:"storage_url"`
	HeaderAttributes dicom.Dataset `json:"header_attributes"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}
