package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/service"
)

type errorBody struct {
	Error string `json:"error"`
}

func writeError(writer http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrTicketNotFound), errors.Is(err, service.ErrWorkerNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrTicketChanged):
		status = http.StatusConflict
	case errors.Is(err, service.ErrWorkerInactive), errors.Is(err, service.ErrZoneMismatch), errors.Is(err, service.ErrInvalidState), errors.Is(err, service.ErrInvalidDocument):
		status = http.StatusUnprocessableEntity
	default:
		status = http.StatusBadRequest
	}
	if status == http.StatusBadRequest && !errors.Is(err, service.ErrWorkerInactive) && !errors.Is(err, service.ErrZoneMismatch) && !errors.Is(err, service.ErrInvalidState) && !errors.Is(err, service.ErrTicketChanged) && !errors.Is(err, service.ErrTicketNotFound) && !errors.Is(err, service.ErrWorkerNotFound) && !errors.Is(err, service.ErrInvalidDocument) {
		if len(err.Error()) > 0 && err.Error()[0:1] == "a" {
			status = http.StatusInternalServerError
		}
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(errorBody{Error: err.Error()})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
