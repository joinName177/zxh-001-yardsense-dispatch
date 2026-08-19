package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func TestCancelledTicketExportStopsBeforeWritingCSV(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	database, err := store.Open(filepath.Join(t.TempDir(), "yardsense.json"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	application := New(database, func() time.Time { return now })
	for index := 0; index < 2; index++ {
		_, err := application.CreateTicket(CreateTicketInput{Title: "Inspect export lane", Description: "Export context cancellation must stop durable yard report processing.", Zone: "north-yard", Priority: model.PriorityNormal, DueAt: now.Add(time.Hour), Operator: "dispatcher", Tags: []string{"report"}})
		if err != nil {
			t.Fatalf("create ticket %d: %v", index, err)
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = application.ExportTicketsCSVContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled export error = %v, want context.Canceled", err)
	}
}
