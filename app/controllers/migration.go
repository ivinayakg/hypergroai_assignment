package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/models"
	"github.com/ivinayakg/hypergroai_assignment/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var AdminSecret = os.Getenv("adminsecret")

type RequestBody struct {
	SecretCode string `json:"secret_code"`
}

func GetMigrations(w http.ResponseWriter, r *http.Request) {
	migrations, err := models.GetMigrations()
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("POST", w, http.StatusCreated)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": migrations})
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check the secret code
	if body.SecretCode != AdminSecret {
		helpers.SendJSONError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	// Parse the multipart form data with a 1 MB file size limit
	err = r.ParseMultipartForm(1 << 20) // 1 MB limit
	if err != nil {
		helpers.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// Upload the file directly to GCS
	err = helpers.UploadCSVFile(w, handler.Filename, file)
	if err != nil {
		helpers.SendJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.SetHeaders("POST", w, http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"messages": "success"})
}

func RunMigrations(w http.ResponseWriter, r *http.Request) {

	// Parse the request body
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Check the secret code
	if body.SecretCode != AdminSecret {
		helpers.SendJSONError(w, http.StatusUnauthorized, "Invalid secret code")
		return
	}

	go utils.RunMigrations()

	helpers.SetHeaders("POST", w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"messages": "success"})
}
