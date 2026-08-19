package service

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func TestImportInvalidSnapshotLeavesStateUnchanged(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	database, err := store.Open(filepath.Join(t.TempDir(), "yardsense.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	application := New(database, func() time.Time { return now })

	worker := model.Worker{ID: "worker-priya", Name: "Priya Rao", Zones: []string{"north-yard"}, Skills: []string{"safety"}, Active: true, Capacity: 2}
	if err := application.AddWorker(worker); err != nil {
		t.Fatalf("add worker: %v", err)
	}
	ticket, err := application.CreateTicket(CreateTicketInput{Title: "Inspect shifted safety rail", Description: "The loading team observed movement at the north loading rail.", Zone: "north-yard", Priority: model.PriorityHigh, DueAt: now.Add(2 * time.Hour), Operator: "dispatcher", Tags: []string{"safety"}})
	if err != nil {
		t.Fatalf("create ticket: %v", err)
	}

	originalSnapshot := database.Snapshot()

	// Snapshot missing version and required collections — must be rejected and must not touch current state.
	bad := store.Document{Version: 0, Tickets: map[string]model.Ticket{}, Workers: map[string]model.Worker{}, Notes: map[string][]model.Note{}}
	if err := application.ImportDocument(bad, "dispatcher"); !errors.Is(err, ErrInvalidDocument) {
		t.Fatalf("import error = %v, want ErrInvalidDocument", err)
	}
	afterSnapshot := database.Snapshot()
	if len(afterSnapshot.Tickets) != 1 || afterSnapshot.Tickets[ticket.ID].Revision != ticket.Revision {
		t.Fatalf("invalid import mutated current state: tickets=%v, want unchanged %v", afterSnapshot.Tickets, originalSnapshot.Tickets)
	}
	if len(afterSnapshot.Workers) != 1 {
		t.Fatalf("invalid import mutated workers: %v", afterSnapshot.Workers)
	}
}
