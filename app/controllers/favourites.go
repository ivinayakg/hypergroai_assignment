package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
	"github.com/ivinayakg/hypergroai_assignment/models"
	"go.mongodb.org/mongo-driver/bson"
)

type AddFavouriteRequestBody struct {
	StockCode string `json:"stockCode"`
}

func AddFavourite(w http.ResponseWriter, r *http.Request) {
	var body AddFavouriteRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.SendJSONError(&w, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value(middleware.UserAuthKey).(*models.User)

	res, err := models.AddFavourite(user.ID.Hex(), body.StockCode)
	if err != nil {
		helpers.SendJSONError(&w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SetHeaders("POST", &w, http.StatusCreated)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": res})
}

func RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stockCode := vars["code"]

	user := r.Context().Value(middleware.UserAuthKey).(*models.User)

	res, err := models.RemoveFavourite(user.ID.Hex(), stockCode)
	if err != nil {
		helpers.SendJSONError(&w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SetHeaders("DELETE", &w, http.StatusAccepted)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": res})
}

func GetUserFavourite(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserAuthKey).(*models.User)

	res, err := models.GetFavouriteStocks(user.ID.Hex())
	if err != nil {
		helpers.SendJSONError(&w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SetHeaders("GET", &w, http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": res})
}
func GetUserFavouriteCodes(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserAuthKey).(*models.User)

	res, err := models.GetFavouriteStocks(user.ID.Hex())
	if err != nil {
		helpers.SendJSONError(&w, http.StatusBadRequest, err.Error())
		return
	}

	results := []string{}
	for _, entry := range res {
		results = append(results, entry.StockCode)
	}

	helpers.SetHeaders("GET", &w, http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": results})
}
