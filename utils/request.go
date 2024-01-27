package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// BindOrReject binds JSON request to provided structure or writes an error response
func BindOrReject(w http.ResponseWriter, r *http.Request, v interface{}, logger *logrus.Logger) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "JSON parse error", logger)
		return err
	}

	// We can further check if the interface has a Validate() method to run validation here...
	validator, ok := v.(interface {
		Validate() error
	})

	if ok {
		err := validator.Validate()
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err, logger)
			return err
		}
	}
	return nil
}
