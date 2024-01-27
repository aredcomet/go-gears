package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type FailedResponse struct {
	Detail string `json:"detail"`
}

func RespondWithError(w http.ResponseWriter, code int, message interface{}, logger *logrus.Logger) {
	w.WriteHeader(code)

	// Check if err is of type string. If so, wrap it with FailedResponse
	if msgStr, ok := message.(string); ok {
		message = FailedResponse{Detail: msgStr}
	}

	if err := json.NewEncoder(w).Encode(message); err != nil {
		logger.Error("Error encoding json: ", err)
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}, logger *logrus.Logger) {
	// Marshal the payload into JSON format, and check error
	response, err := json.Marshal(payload)
	if err != nil {
		// Log and respond with an error if JSON marshalling fails
		logger.Error("Error encoding json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
