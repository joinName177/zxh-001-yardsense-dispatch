package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/service"
)

type ticketInput struct {
	Title, Description, Zone, Priority, Operator, DueAt string
	Tags                                                []string `json:"tags"`
}
type assignmentInput struct {
	WorkerID string `json:"worker_id"`
	Operator string `json:"operator"`
	Revision int64  `json:"revision"`
}
type statusInput struct {
	Status   string `json:"status"`
	Operator string `json:"operator"`
	Revision int64  `json:"revision"`
}
type noteInput struct {
	Author   string `json:"author"`
	Body     string `json:"body"`
	Internal bool   `json:"internal"`
}

func decode(request *http.Request, destination any) error {
	return json.NewDecoder(http.MaxBytesReader(nil, request.Body, 1<<20)).Decode(destination)
}

func (s *Server) listTickets(writer http.ResponseWriter, request *http.Request) {
	filter := model.TicketFilter{Zone: request.URL.Query().Get("zone"), AssigneeID: request.URL.Query().Get("assignee"), Tag: request.URL.Query().Get("tag"), Overdue: request.URL.Query().Get("overdue") == "true"}
	if value := request.URL.Query().Get("status"); value != "" {
		filter.Status = model.Status(value)
	}
	writeJSON(writer, http.StatusOK, s.service.Tickets(filter))
}

func (s *Server) createTicket(writer http.ResponseWriter, request *http.Request) {
	var body ticketInput
	if err := decode(request, &body); err != nil {
		writeError(writer, err)
		return
	}
	dueAt, err := time.Parse(time.RFC3339, body.DueAt)
	if err != nil {
		writeError(writer, err)
		return
	}
	ticket, err := s.service.CreateTicket(service.CreateTicketInput{Title: body.Title, Description: body.Description, Zone: body.Zone, Operator: body.Operator, Priority: model.Priority(body.Priority), DueAt: dueAt, Tags: body.Tags})
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, ticket)
}

func (s *Server) getTicket(writer http.ResponseWriter, request *http.Request) {
	ticket, err := s.service.Ticket(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, ticket)
}

func (s *Server) assignTicket(writer http.ResponseWriter, request *http.Request) {
	var body assignmentInput
	if err := decode(request, &body); err != nil {
		writeError(writer, err)
		return
	}
	ticket, err := s.service.Assign(model.Assignment{TicketID: request.PathValue("id"), WorkerID: body.WorkerID, AssignedBy: body.Operator, ExpectedRev: body.Revision, AssignedAt: time.Now().UTC()})
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, ticket)
}

func (s *Server) changeStatus(writer http.ResponseWriter, request *http.Request) {
	var body statusInput
	if err := decode(request, &body); err != nil {
		writeError(writer, err)
		return
	}
	ticket, err := s.service.ChangeStatus(service.ChangeStatusInput{TicketID: request.PathValue("id"), Operator: body.Operator, ExpectedRevision: body.Revision, Status: model.Status(body.Status)})
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, ticket)
}

func (s *Server) notes(writer http.ResponseWriter, request *http.Request) {
	notes, err := s.service.Notes(request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, notes)
}

func (s *Server) addNote(writer http.ResponseWriter, request *http.Request) {
	var body noteInput
	if err := decode(request, &body); err != nil {
		writeError(writer, err)
		return
	}
	note := service.NewNote("note-"+time.Now().UTC().Format("20060102150405.000000000"), request.PathValue("id"), body.Author, body.Body, body.Internal, time.Now())
	if err := s.service.AddNote(note); err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, note)
}
