package db

import (
	"fmt"
	"time"

	"github.com/suyashkumar/dicom"

	"pockethealth/dicom/model"
)

const StorageLocation string = "images"

// I'm assuming that all the user data is going to be coming from a separate microservice,
// so I'm storing the user ids but I don't expect that this service will have much more
// information about them. That will need to be fetched from somewhere else.

// Not implementing this, but writing out probably what it should include.
func CreateDicomImage(imageId string, uploaderId string, patientId string, headerAttributes []*dicom.Element) model.DicomImage {
	// I'm imagining that the table for this is going to be something like
	// - id (uuid): Primary Key use the imageId
	// - uploaded_by_user_id (uuid): The authed user that uploaded the image: uploaderId
	// - patient_id (uuid): Id of the patient that the image is of: patientID
	// - storage_url (string): Wherever we are saving these images
	// - header_attributes: jsonb column with all the extracted data
	//   - if this needs to get queried often then we should move its own table, or columns
	// - created_at (timestamp)
	// - updated_at (timestamp)

	return model.DicomImage{
		Id:               imageId,
		UploadedByUserId: uploaderId,
		PatientId:        patientId,
		StorageUrl:       fmt.Sprintf("%s/%s.dcm", StorageLocation, imageId),
		HeaderAttributes: headerAttributes,
		CreatedAt:        time.Now(),
	}
}

// Again not implementing this
func CreatePngImage(imageId string, dicomImageId string, patientId string) model.PngImage {
	// Much the same as the DicomImage, but as no one actually uploads this I'm omitting that.
	// - id (uuid): Primary Key use the imageId
	// - dicom_image_id (uuid): Foreign key to the DicomImage table
	// - patient_id (uuid): Id of the patient that the image is of: patientID
	// - storage_url (string): Wherever we are saving these images
	// - created_at (timestamp)
	// - updated_at (timestamp)

	// Originally I was going to use a shared primary key between this and the DicomImage, but
	// after I realised that the pixel data can return an array of frames I thought that it
	// might be better to make it have it's own id, so that you can have many pngs from one dicom.
	// Maybe an incorrect assumption, but I don't really know enough about these images to know
	// any better.

	return model.PngImage{
		Id:           imageId,
		DicomImageId: dicomImageId,
		PatientId:    patientId,
		StorageUrl:   fmt.Sprintf("%s/%s.png", StorageLocation, imageId),
		CreatedAt:    time.Now(),
	}
}

func GetPatientIdForPngImage(imageId string) string {
	// If the DB was in place then you just do the query and return the id.
	// Because I haven't done any of that I'm just returning the image id back to you.
	return imageId
}
