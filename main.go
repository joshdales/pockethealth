package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func main() {
	router := chi.NewRouter()

	router.Use(checkUserAuthentication)

	// Upload a DICOM image
	router.Post("/image", handleDicomImageUpload)
	// Return the image as a PNG
	router.Get("/image/{imageId}", handleGetImageById)

	err := http.ListenAndServe(":3333", router)
	if err != nil {
		fmt.Printf("Error while opening the server: %v", err)
	}
}

// I'm not actually implementing this, I'm assuming that a different
// microservice should be responsible for this.
func checkUserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call whatever service we need to make sure that you have access
		// Store the authed user's id in context
		ctx := context.WithValue(r.Context(), "userId", uuid.New())
		r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// Again I'm assuming that this is going to be handled by a separate
// microservice and I'm not actually going to implement it.
func userHasAccessToPatient(authedUserId string, permittedRoles []string, patientId string) bool {
	// Check that the user has access to the patient, and the role for that endpoint.
	// - If they are a patient then they must either be that patient or their guardian
	// - If they are a clinician then they should have access to that patient,
	//   i.e if you are a clinician who is not involved with that patient's care then you
	//   shouldn't have access to their medical records.
	// If we pass the check then great we can complete the request.
	return true
}

func handleDicomImageUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(512)
	if err != nil {
		http.Error(w, "Could not parse the form data", http.StatusBadRequest)
	}

	patientId := r.Form.Get("patientId")
	userId := r.Context().Value("userId").(string)
	allowedRoles := make([]string, 1)
	allowedRoles[0] = "clinician"
	if !userHasAccessToPatient(userId, allowedRoles, patientId) {
		http.Error(w, "You do not have access to this", http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageId := uuid.New()
	filePath := fmt.Sprintf("images/%s.dcm", imageId)

	err = os.Mkdir("images", 0750)
	if err != nil && !os.IsExist(err) {
		http.Error(w, "Uh oh! Something has gone wrong", http.StatusInternalServerError)
	}
	err = os.WriteFile(filePath, file)
	if err != nil {
		http.Error(w, "Uh oh! Something has gone wrong", http.StatusInternalServerError)
	}
}

func handleGetImageById(w http.ResponseWriter, r *http.Request) {
	imageId := chi.URLParam(r, "imageId")
	patientId := getPatientIdForImage(imageId)
	userId := r.Context().Value("userId").(string)
	allowedRoles := make([]string, 2)
	allowedRoles[0] = "clinician"
	allowedRoles[1] = "patient"
	if !userHasAccessToPatient(userId, allowedRoles, patientId) {
		http.Error(w, "You do not have access to this", http.StatusForbidden)
		return
	}
	imagePath := fmt.Sprintf("images/%s.png", imageId)
	image, err := os.ReadFile(imagePath)
	if err != nil {
		if !os.IsExist(err) {
			http.Error(w, "Requested image does not exist", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(image)
}

func getPatientIdForImage(imageId string) string {
	// TODO: implement some sort of DB to hold this information
	// for now I'm just using the image id as the patient id
	patientId := imageId
	return patientId
}
