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

type UserRole = string

const (
	Patient   UserRole = "patient"
	Clinician UserRole = "clinician"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(checkUserAuthentication)

	// Upload a DICOM image
	router.Post("/dicom", handleDicomImageUpload)
	// Return information about the DICOM image
	router.Get("/dicom/{imageId}/header_attributes", handleGetDicomImageById)
	// Return PNG version of the image
	router.Get("/png/{imageId}/image", handleGetPngImageById)

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
func userHasAccessToPatient(authedUserId string, permittedRoles []UserRole, patientId string) bool {
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

	patientId := r.FormValue("patient_id")
	userId := r.Context().Value("userId").(string)
	allowedRoles := []UserRole{Clinician}
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

	img := db.CreateDicomImage(imgId, userId, patientId, nil)

	dataset, err := dicom.ParseFile(fmt.Sprintf("%s/%s.dcm", db.StorageLocation, imgId), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	img.HeaderAttributes = dataset.Elements

	err = util.ConvertDicomToPngAndUpload(&dataset, imgId, patientId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(img)
}

func handleGetDicomImageById(w http.ResponseWriter, r *http.Request) {
	imageId := chi.URLParam(r, "imageId")
	patientId := db.GetPatientIdForPngImage(imageId)
	userId := r.Context().Value("userId").(string)
	allowedRoles := []UserRole{Clinician}
	if !userHasAccessToPatient(userId, allowedRoles, patientId) {
		http.Error(w, "You do not have access to this", http.StatusForbidden)
		return
	}

	// If you had a DB you could just read that but as we don't I have to re-parse the file,
	// and then I'm using the create function just to make life a little simpler.
	imagePath := fmt.Sprintf("%s/%s.dcm", db.StorageLocation, imageId)
	dataset, err := dicom.ParseFile(imagePath, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var filteredDataset []*dicom.Element
	// If no query is provided then return all the elements
	if len(r.URL.Query()) < 1 {
		filteredDataset = dataset.Elements
	} else {
		for query := range r.URL.Query() {
			for _, element := range dataset.Elements {
				if fmt.Sprintf("%s", element.Tag) == query {
					filteredDataset = append(filteredDataset, element)
				}
			}
		}
	}

	img := db.CreateDicomImage(imageId, userId, patientId, filteredDataset)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(img.HeaderAttributes)
}

func handleGetPngImageById(w http.ResponseWriter, r *http.Request) {
	imageId := chi.URLParam(r, "imageId")
	patientId := db.GetPatientIdForPngImage(imageId)
	userId := r.Context().Value("userId").(string)
	allowedRoles := []UserRole{Clinician, Patient}
	if !userHasAccessToPatient(userId, allowedRoles, patientId) {
		http.Error(w, "You do not have access to this", http.StatusForbidden)
		return
	}

	imagePath := fmt.Sprintf("%s/%s.png", db.StorageLocation, imageId)
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
