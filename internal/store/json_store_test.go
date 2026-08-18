package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func TestOpenRestoresPersistedWorkerAndTicket(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	database, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Date(2026, 8, 18, 8, 0, 0, 0, time.UTC)
	worker := model.Worker{ID: "worker-1", Name: "Shift Worker", Zones: []string{"zone-1"}, Skills: []string{"safety"}, Active: true, Capacity: 1}
	if err := database.UpsertWorker(worker); err != nil {
		t.Fatalf("save worker: %v", err)
	}
	ticket := model.Ticket{ID: "ticket-1", Title: "Inspect loading dock", Description: "Loading dock inspection requested after a near miss event.", Zone: "zone-1", Priority: model.PriorityHigh, Status: model.StatusOpen, DueAt: now.Add(time.Hour), CreatedAt: now, UpdatedAt: now, Revision: 1}
	activity := model.Activity{TicketID: ticket.ID, Kind: model.ActivityCreated, Actor: "operator", Message: "ticket created", At: now}
	if err := database.CreateTicket(ticket, activity); err != nil {
		t.Fatalf("save ticket: %v", err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	loaded, err := reopened.ReadTicket(ticket.ID)
	if err != nil {
		t.Fatalf("read reopened ticket: %v", err)
	}
	if loaded.Title != ticket.Title || loaded.Revision != ticket.Revision {
		t.Fatalf("reopened ticket = %#v, want persisted ticket", loaded)
	}
}
