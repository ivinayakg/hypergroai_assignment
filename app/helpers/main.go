package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var methodChoices = map[string]string{
	"get":   "GET",
	"post":  "POST",
	"patch": "PATCH",
	"del":   "DELETE",
}

func SetHeaders(type_ string, w *http.ResponseWriter, status int) {
	method := methodChoices[type_]
	if method == "" {
		method = "GET"
	}

	(*w).WriteHeader(status)
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	if method != "GET" {
		(*w).Header().Set("Access-Control-Allow-Methods", method)
		(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func SendJSONError(w *http.ResponseWriter, statusCode int, errorMessage string) {
	errorResponse := ErrorResponse{Error: errorMessage}
	(*w).Header().Set("content-Type", "application/json")
	(*w).WriteHeader(statusCode)
	json.NewEncoder((*w)).Encode(errorResponse)
}

func ContainsString(arr *[]string, target *string) bool {
	for _, s := range *arr {
		if strings.Contains(s, *target) {
			return true
		}
	}
	return false
}

func GetUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	// If X-Forwarded-For header is empty (not behind a proxy), get the IP from RemoteAddr
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

func ConvertFilenameToDate(filename string) (*time.Time, error) {
	// Extract year, month, and day from the filename
	dayStr := filename[2:4]
	monthStr := filename[4:6]
	yearStr := filename[6:8]

	// Convert strings to integers
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return nil, fmt.Errorf("error converting year: %w", err)
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return nil, fmt.Errorf("error converting month: %w", err)
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return nil, fmt.Errorf("error converting day: %w", err)
	}

	// Construct a time.Time object
	date := time.Date(2000+year, time.Month(month), day, 0, 0, 0, 0, time.FixedZone("UTC+5:30", 5*60*60+30*60))
	return &date, nil
}

func ConvertTimeIntoDate(t *time.Time) (string, error) {
	// Format the date as "02/01/2006"
	formattedDate := t.Format("02/01/2006")

	return formattedDate, nil
}

func ConvertStringToFloat(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}

func ConvertStringToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return 0
	}

	return result
}
