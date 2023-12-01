package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"

	"pockethealth/dicom/db"
	"pockethealth/dicom/util"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
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
		userId := uuid.New().String()
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
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
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		http.Error(w, "Could not parse the form data", http.StatusBadRequest)
	}

	patientId := r.FormValue("patientId")
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

	imgId, err := util.UploadImage(file, patientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img := db.CreateDicomImage(imgId, userId, patientId)

	dataset, err := dicom.Parse(file, 2<<20, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img = db.UpdateDicomImage(img, dataset)

	err = util.ConvertDicomToPngAndUpload(&dataset, imgId, patientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(img)
}

func handleGetImageById(w http.ResponseWriter, r *http.Request) {
	imageId := chi.URLParam(r, "imageId")
	patientId := db.GetPatientIdForPngImage(imageId)
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
		if os.IsNotExist(err) {
			http.Error(w, "Requested image does not exist", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(image)
}
