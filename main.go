package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := http.NewServeMux()
	// Upload a DICOM image
	server.HandleFunc("/patient/:patientId/image", handleDicomUpload)
	// Query DICOM header attribute
	server.HandleFunc("/patient/:patientId/image/:imageId", handleDicomAttributeQuery)
	// Return the converted image as a PNG
	server.HandleFunc("/patient/:patientId/image/:imageId/png", handlePngImage)

	err := http.ListenAndServe(":3333", server)
	if err != nil {
		fmt.Printf("Error while opening the server: %v", err)
	}
}

// I'm not actually implementing this, just writing out what we should be checking
func authedUserHasAccessToPatient() bool {
	// - Make sure that the user is authenticated
	// - Then get the patient ID from the url
	// - Check that the patient exists
	// - Check that the logged in user has access to the patient:
	//   - If they are a patient then they must either be that patient or their guardian
	//   - If they are a clinician then they should have access to that patient,
	//     i.e if you are a clinician who is not involved with that patient's care then you
	//     shouldn't have access to their medical records.
	// - If we pass all these checks then great we can complete the request.
	return true
}

func handleDicomUpload(w http.ResponseWriter, r *http.Request) {
	// Make sure that the user has access to this image, they should be a clinician that has access to the patient.
	if !authedUserHasAccessToPatient() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}
}

func handleDicomAttributeQuery(w http.ResponseWriter, r *http.Request) {
	// Make sure that the user has access to this image, they should be a clinician that has access to the patient.
	if !authedUserHasAccessToPatient() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}
}

func handlePngImage(w http.ResponseWriter, r *http.Request) {
	// Only the patient or clinician with access to the patient should have access this image.
	if !authedUserHasAccessToPatient() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	} else if r.Method != http.MethodGet {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}
}
