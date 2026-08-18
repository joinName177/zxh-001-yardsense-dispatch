package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func TestListedTicketTagsCannotMutatePersistedTicket(t *testing.T) {
	path := filepath.Join(t.TempDir(), "yardsense.json")
	database, err := Open(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	now := time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)
	ticket := model.Ticket{
		ID: "ticket-owned-tags", Title: "Inspect spill kit cabinet", Description: "Night shift reported a damaged spill kit cabinet seal.", Zone: "north-yard",
		Priority: model.PriorityHigh, Status: model.StatusOpen, DueAt: now.Add(time.Hour), CreatedAt: now, UpdatedAt: now, Revision: 1, Tags: []string{"safety"},
	}
	activity := model.Activity{TicketID: ticket.ID, Kind: model.ActivityCreated, Actor: "dispatcher", Message: "ticket created", At: now}
	if err := database.CreateTicket(ticket, activity); err != nil {
		t.Fatalf("create ticket: %v", err)
	}
	listed := database.ListTickets(model.TicketFilter{}, now)
	if len(listed) != 1 {
		t.Fatalf("listed tickets = %d, want 1", len(listed))
	}
	listed[0].Tags[0] = "electrical"
	worker := model.Worker{ID: "worker-persist", Name: "Persist Worker", Zones: []string{"north-yard"}, Skills: []string{"safety"}, Active: true, Capacity: 2}
	if err := database.UpsertWorker(worker); err != nil {
		t.Fatalf("persist unrelated worker: %v", err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	persisted, err := reopened.ReadTicket(ticket.ID)
	if err != nil {
		t.Fatalf("read ticket after reopen: %v", err)
	}
	if persisted.Tags[0] != "safety" {
		t.Fatalf("listed ticket mutated persisted tag = %q, want %q", persisted.Tags[0], "safety")
	}
}
