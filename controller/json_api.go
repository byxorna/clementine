package controller

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	OK   bool        `json:"ok"`
	Data interface{} `json:"data"`
}

func jsonError(w http.ResponseWriter, errorString string) {
	jsonErrorStatus(w, errorString, http.StatusBadRequest)
}

func jsonErrorStatus(w http.ResponseWriter, errorString string, status int) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	d := JsonResponse{false, errorString}
	js, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(js)
}

func jsonSuccess(w http.ResponseWriter, data interface{}) {
	jsonSuccessStatus(w, data, http.StatusOK)
}

func jsonSuccessStatus(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(http.StatusOK)
	d := JsonResponse{true, data}
	js, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(js)
}
