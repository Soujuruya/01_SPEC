package httphelper

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	resp := APIResponse{
		Status: "ok",
		Data:   data,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func WriteError(w http.ResponseWriter, err error, statusCode int) {
	resp := APIResponse{
		Status: "error",
		Error:  err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}
