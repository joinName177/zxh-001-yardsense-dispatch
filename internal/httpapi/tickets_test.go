package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/service"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func TestAssignReturnsConflictWhenTicketWasUpdated(t *testing.T) {
	now := time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)
	database, err := store.Open(filepath.Join(t.TempDir(), "yardsense.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	application := service.New(database, func() time.Time { return now })
	worker := model.Worker{ID: "worker-maya", Name: "Maya Chen", Zones: []string{"north-yard"}, Skills: []string{"safety"}, Active: true, Capacity: 2}
	if err := application.AddWorker(worker); err != nil {
		t.Fatalf("add worker: %v", err)
	}
	ticket, err := application.CreateTicket(service.CreateTicketInput{Title: "Inspect shifted safety rail", Description: "The loading team observed movement at the north loading rail.", Zone: "north-yard", Priority: model.PriorityHigh, DueAt: now.Add(2 * time.Hour), Operator: "dispatcher", Tags: []string{"safety"}})
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}
	if _, err := application.Assign(model.Assignment{TicketID: ticket.ID, WorkerID: worker.ID, AssignedBy: "dispatcher", AssignedAt: now, ExpectedRev: ticket.Revision}); err != nil {
		t.Fatalf("initial assignment: %v", err)
	}
	body, err := json.Marshal(map[string]any{"worker_id": worker.ID, "operator": "dispatcher", "revision": ticket.Revision})
	if err != nil {
		t.Fatalf("encode body: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/tickets/"+ticket.ID+"/assign", bytes.NewReader(body))
	response := httptest.NewRecorder()
	New(application).ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("stale dispatcher update returned %d, want %d; body=%s", response.Code, http.StatusConflict, response.Body.String())
	}
}

func TestImportInvalidSnapshotReturnsUnprocessableEntity(t *testing.T) {
	database, err := store.Open(filepath.Join(t.TempDir(), "yardsense.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	application := service.New(database, time.Now)
	request := httptest.NewRequest(http.MethodPost, "/api/import/document", bytes.NewBufferString(`{"version":1,"tickets":null,"workers":{},"notes":{}}`))
	request.Header.Set("X-Operator", "dispatcher")
	response := httptest.NewRecorder()
	New(application).ServeHTTP(response, request)
	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid snapshot returned %d, want %d; body=%s", response.Code, http.StatusUnprocessableEntity, response.Body.String())
	}
}
