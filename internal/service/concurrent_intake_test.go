package service

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func TestConcurrentTicketIntakeKeepsBothTickets(t *testing.T) {
	now := time.Date(2026, 8, 19, 11, 0, 0, 0, time.UTC)
	database, err := store.Open(filepath.Join(t.TempDir(), "yardsense.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	application := New(database, func() time.Time { return now })
	start := make(chan struct{})
	errors := make(chan error, 2)
	var group sync.WaitGroup
	for _, title := range []string{"Inspect north loading ramp", "Inspect south loading ramp"} {
		group.Add(1)
		go func(title string) {
			defer group.Done()
			<-start
			_, err := application.CreateTicket(CreateTicketInput{Title: title, Description: "Two dispatchers submitted independent safety inspections at the same instant.", Zone: "north-yard", Priority: model.PriorityHigh, DueAt: now.Add(time.Hour), Operator: "dispatcher", Tags: []string{"safety"}})
			errors <- err
		}(title)
	}
	close(start)
	group.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatalf("concurrent intake returned error: %v", err)
		}
	}
	if tickets := application.Tickets(model.TicketFilter{}); len(tickets) != 2 {
		t.Fatalf("persisted tickets = %d, want 2", len(tickets))
	}
}
