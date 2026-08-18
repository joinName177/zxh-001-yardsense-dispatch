package httpapi

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/service"
)

//go:embed static/*
var staticFiles embed.FS

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(application *service.Service) *Server {
	server := &Server{service: application, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /api/dashboard", s.dashboard)
	s.mux.HandleFunc("GET /api/triage", s.triageQueue)
	s.mux.HandleFunc("GET /api/reports/operations", s.operationsReport)
	s.mux.HandleFunc("GET /api/reports/workload", s.workloadReport)
	s.mux.HandleFunc("GET /api/export/tickets.csv", s.exportCSV)
	s.mux.HandleFunc("GET /api/export/document", s.exportDocument)
	s.mux.HandleFunc("POST /api/import/document", s.importDocument)
	s.mux.HandleFunc("POST /api/demo/seed", s.seedDemo)
	s.mux.HandleFunc("GET /api/tickets", s.listTickets)
	s.mux.HandleFunc("POST /api/tickets", s.createTicket)
	s.mux.HandleFunc("GET /api/tickets/{id}", s.getTicket)
	s.mux.HandleFunc("POST /api/tickets/{id}/assign", s.assignTicket)
	s.mux.HandleFunc("POST /api/tickets/{id}/status", s.changeStatus)
	s.mux.HandleFunc("GET /api/tickets/{id}/timeline", s.timeline)
	s.mux.HandleFunc("GET /api/tickets/{id}/recommendations", s.recommendations)
	s.mux.HandleFunc("GET /api/tickets/{id}/notes", s.notes)
	s.mux.HandleFunc("POST /api/tickets/{id}/notes", s.addNote)
	s.mux.HandleFunc("GET /api/workers", s.listWorkers)
	s.mux.HandleFunc("POST /api/workers", s.addWorker)
	content, _ := fs.Sub(staticFiles, "static")
	s.mux.Handle("GET /", http.FileServerFS(content))
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.URL.Path, "/api/") {
		writer.Header().Set("Cache-Control", "no-store")
	}
	securityHeaders(accessLog(s.mux)).ServeHTTP(writer, request)
}

func (s *Server) health(writer http.ResponseWriter, request *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}
