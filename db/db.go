package db

// I'm assuming that all the user data is going to be coming from a separate microservice,
// so I'm storing the user ids but I don't expect that this service will have much more
// information about them. That will need to be fetched from somewhere else.

// Not implementing this, but writing out probably what it should include.
func CreateDicomImage(imageId string, uploaderId string, patientId string) {
	// I'm imagining that the table for this is going to be something like
	// - id (uuid): Primary Key use the imageId
	// - uploaded_by (uuid): The authed user that uploaded the image: uploaderId
	// - patient_id (uuid): Id of the patient that the image is of: patientID
	// - storage_url (string): Wherever we are saving these images
	// - header_attributes: jsonb column with all the extracted data
	//   - if this needs to get queried often then we should move its own table, or columns
	// - created_at (timestamp)
	// - updated_at (timestamp)
}

func UpdateDicomImage(imageId string, dataToUpdate any) {
	// This would just be an update function, but I'm not creating all the model and types
}

// Again not implementing this
func CreatePngImage(imageId string, dicomImageId string, patientId string) {
	// Much the same as the DicomImage
	// - id (uuid): Primary Key use the imageId
	// - dicom_image_id (uuid): Foreign key to the DicomImage table
	// - patient_id (uuid): Id of the patient that the image is of: patientID
	// - storage_url (string): Wherever we are saving these images
	// - header_attributes: jsonb column with all the extracted data
	// - created_at (timestamp)
	// - updated_at (timestamp)
}

func GetPatientIdForPngImage(imageId string) string {
	// If the DB was in place then you just do the query and return the id.
	// Because I haven't done any of that I'm just returning the image id back to you.
	return imageId
}
