package main

import (
	"context"
	"fmt"
	"image/png"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
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

	var fileData []byte
	file.Read(fileData)

	defer file.Close()

	imageId := uuid.New()
	filePath := fmt.Sprintf("images/%s.dcm", imageId)

	err = os.Mkdir("images", 0750)
	if err != nil && !os.IsExist(err) {
		http.Error(w, "Couldn't create images folder", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filePath, fileData, 0660)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	// I'm sure that you can probably do this with the file stream but I don't know Go well enough to do it.
	// So I'm using the example from https://pkg.go.dev/github.com/suyashkumar/dicom@v1.0.7#example-package-ReadFile
	// as a base for how to do this.
	dataset, err := dicom.ParseFile(filePath, nil)
	if err != nil {
		http.Error(w, "Could not find Image to parse", http.StatusInternalServerError)
		return
	}
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		http.Error(w, "Could not get pixel data for image", http.StatusInternalServerError)
		return
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for _, frame := range pixelDataInfo.Frames {
		img, err := frame.GetImage()
		if err != nil {
			http.Error(w, "Could not get image frame", http.StatusInternalServerError)
			return
		}
		imgFile, err := os.Create(fmt.Sprintf("images/%s.png", uuid.New()))
		if err != nil {
			http.Error(w, "Could not create png image file", http.StatusInternalServerError)
			return
		}
		err = png.Encode(imgFile, img)
		if err != nil {
			http.Error(w, "Could not encode png image", http.StatusInternalServerError)
			return
		}
		err = imgFile.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func getPatientIdForImage(imageId string) string {
	// TODO: implement some sort of DB to hold this information
	// for now I'm just using the image id as the patient id
	patientId := imageId
	return patientId
}
