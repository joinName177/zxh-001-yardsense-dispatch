package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func (s *Server) operationsReport(writer http.ResponseWriter, request *http.Request) {
	hours := 24
	if value := request.URL.Query().Get("hours"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			writeError(writer, err)
			return
		}
		hours = parsed
	}
	report, err := s.service.OperationsReport(hours)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, report)
}

func (s *Server) workloadReport(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, s.service.WorkloadReport())
}

func (s *Server) timeline(writer http.ResponseWriter, request *http.Request) {
	timeline, err := s.service.TicketTimeline(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, timeline)
}

func (s *Server) exportCSV(writer http.ResponseWriter, request *http.Request) {
	contents, err := s.service.ExportTicketsCSVContext(request.Context())
	if err != nil {
		writeError(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "text/csv; charset=utf-8")
	writer.Header().Set("Content-Disposition", "attachment; filename=yardsense-tickets.csv")
	_, _ = writer.Write(contents)
}

func (s *Server) exportDocument(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, s.service.ExportDocument())
}

func (s *Server) importDocument(writer http.ResponseWriter, request *http.Request) {
	operator := request.Header.Get("X-Operator")
	var document store.Document
	if err := json.NewDecoder(http.MaxBytesReader(writer, request.Body, 8<<20)).Decode(&document); err != nil {
		writeError(writer, err)
		return
	}
	if err := s.service.ImportDocument(document, operator); err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusNoContent, nil)
}

func (s *Server) seedDemo(writer http.ResponseWriter, request *http.Request) {
	result, err := s.service.SeedDemo(request.Header.Get("X-Operator"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, result)
}
