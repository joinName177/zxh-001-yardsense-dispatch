package httpapi

import "net/http"

func (s *Server) triageQueue(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]any{"summary": s.service.TriageSummary(), "items": s.service.TriageQueue()})
}

func (s *Server) recommendations(writer http.ResponseWriter, request *http.Request) {
	items, err := s.service.RecommendAssignment(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, items)
}
