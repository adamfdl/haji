package main

import (
	"encoding/json"
	"net/http"
)

func writeJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func badRequestResponse(w http.ResponseWriter) {
	http.Error(w, `{"error": "bad request"}`, http.StatusBadRequest)
}

func internalServerErrorResponse(w http.ResponseWriter) {
	http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
}
