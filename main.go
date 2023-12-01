package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func main() {
	router := chi.NewRouter()

	router.Use(checkUserAuthentication)

	// Upload a DICOM image
	router.Post("/image", handleDicomImageUpload)
	// Return the image
	router.Get("/image/{imageId}", handlePngImage)

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
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
}

func handlePngImage(w http.ResponseWriter, r *http.Request) {}
