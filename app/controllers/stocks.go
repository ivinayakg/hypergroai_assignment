package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/models"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type RequestBodyStocks struct {
	Date string `json:"date"`
}

func GetStocksUnverified(w http.ResponseWriter, r *http.Request) {
	migration, err := models.GetLastMigration()
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	stocks, err := models.GetStocksData(1, migration.DataDate, "", int64(10), false)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("GET", w, http.StatusOK)

	if len(stocks) == 0 {
		json.NewEncoder(w).Encode(bson.M{"message": "success", "data": "[]"})
		return
	}

	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": stocks})
}

func GetStocks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page := helpers.ConvertStringToInt(queryParams.Get("page"))
	size := helpers.ConvertStringToInt(queryParams.Get("size"))
	search := queryParams.Get("s")
	date := queryParams.Get("date")

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	// Parse the string into a time.Time value
	if date == "" {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("date is a required query_param").Error())
		return
	}

	// get migration
	migration, err := models.GetMigrationForTime(date)
	if migration == nil {
		fmt.Println(err)
		migration, err = models.GetLastMigration()
		if err != nil {
			helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	stocks, err := models.GetStocksData(page, migration.DataDate, search, int64(size), false)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("GET", w, http.StatusOK)

	if len(stocks) == 0 {
		json.NewEncoder(w).Encode(bson.M{"message": "success", "data": []string{}})
		return
	}

	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": stocks, "pagination": bson.M{"page": page, "size": 15, "next": page + 1}})
}

func GetTopStocks(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	page := helpers.ConvertStringToInt(queryParams.Get("page"))
	size := helpers.ConvertStringToInt(queryParams.Get("size"))
	search := queryParams.Get("s")
	date := queryParams.Get("date")

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	// Parse the string into a time.Time value
	if date == "" {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("date is a required query_param").Error())
		return
	}

	// get migration
	migration, err := models.GetMigrationForTime(date)
	if migration == nil {
		fmt.Println(err)
		migration, err = models.GetLastMigration()
		if err != nil {
			helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	stocks, err := models.GetStocksData(page, migration.DataDate, search, int64(size), true)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("GET", w, http.StatusOK)

	if len(stocks) == 0 {
		json.NewEncoder(w).Encode(bson.M{"message": "success", "data": []string{}})
		return
	}

	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": stocks, "pagination": bson.M{"page": page, "size": 15, "next": page + 1}})
}

func GetStockDetail(w http.ResponseWriter, r *http.Request) {
	stockCode, found := mux.Vars(r)["id"]
	if !found {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("stock Code is required in the url").Error())
		return
	}

	stock, err := models.GetStockDetail(stockCode)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("GET", w, http.StatusOK)
	json.NewEncoder(w).Encode(bson.M{"message": "success", "data": stock})
}
