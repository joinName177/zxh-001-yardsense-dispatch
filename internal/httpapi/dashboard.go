package httpapi

import "net/http"

func (s *Server) dashboard(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, s.service.Dashboard())
}
