package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func (s *Server) listWorkers(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, s.service.Workers())
}

func (s *Server) addWorker(writer http.ResponseWriter, request *http.Request) {
	var worker model.Worker
	if err := json.NewDecoder(request.Body).Decode(&worker); err != nil {
		writeError(writer, err)
		return
	}
	if err := s.service.AddWorker(worker); err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, worker)
}
