package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Route("/patient/{patientId}", func(r chi.Router) {
		r.Use(authedUserHasAccessToPatient)

		// Upload a DICOM image
		r.Post("/image", handleDicomUpload)
		// Query DICOM header attribute
		r.Get("/patient/:patientId/image/:imageId", handleDicomAttributeQuery)
		// Return the converted image as a PNG
		r.Get("/patient/:patientId/image/:imageId/png", handlePngImage)
	})

	err := http.ListenAndServe(":3333", router)
	if err != nil {
		fmt.Printf("Error while opening the server: %v", err)
	}
}

// I'm not actually implementing this, just writing out what we should be checking.
// I'm assuming that a different microservice should be responsible for these checks.
func authedUserHasAccessToPatient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure that the user is authenticated
		// Then get the patient ID from the url and check that the patient exists

		// Check that the logged in user has access to the patient:
		// - If they are a patient then they must either be that patient or their guardian
		// - If they are a clinician then they should have access to that patient,
		//   i.e if you are a clinician who is not involved with that patient's care then you
		//   shouldn't have access to their medical records.

		// If we pass all these checks then great we can complete the request.
		next.ServeHTTP(w, r)
	})
}

func handleDicomUpload(w http.ResponseWriter, r *http.Request) {}

func handleDicomAttributeQuery(w http.ResponseWriter, r *http.Request) {}

func handlePngImage(w http.ResponseWriter, r *http.Request) {}
