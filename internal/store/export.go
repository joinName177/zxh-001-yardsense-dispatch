package store

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func (s *JSONStore) ExportTicketsCSV(writer io.Writer, now time.Time) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "title", "zone", "priority", "status", "assignee", "due_at", "revision", "overdue", "tags"}); err != nil {
		return fmt.Errorf("write csv header: %w", err)
	}
	for _, ticket := range s.ListTickets(model.TicketFilter{}, now) {
		record := []string{ticket.ID, ticket.Title, ticket.Zone, string(ticket.Priority), string(ticket.Status), ticket.AssigneeID, ticket.DueAt.Format(time.RFC3339), strconv.FormatInt(ticket.Revision, 10), strconv.FormatBool(ticket.IsOverdue(now)), strings.Join(ticket.Tags, ",")}
		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("write ticket %s: %w", ticket.ID, err)
		}
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("flush csv: %w", err)
	}
	return nil
}

type ExportDocument struct {
	ExportedAt time.Time `json:"exported_at"`
	Document   Document  `json:"document"`
}

func (s *JSONStore) Export() ExportDocument {
	return ExportDocument{ExportedAt: time.Now().UTC(), Document: s.Snapshot()}
}

func (s *JSONStore) Import(document Document) error {
	if err := validateDocument(document); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	original := s.doc
	s.doc = document.Clone()
	if err := s.persistLocked(); err != nil {
		s.doc = original
		return err
	}
	return nil
}

func validateDocument(document Document) error {
	if document.Version != 1 || document.Tickets == nil || document.Workers == nil || document.Notes == nil {
		return ErrInvalidDocument
	}
	for id, ticket := range document.Tickets {
		if id != ticket.ID || ticket.Revision < 1 || !ticket.Status.Valid() || !ticket.Priority.Valid() {
			return fmt.Errorf("%w: ticket %q", ErrInvalidDocument, id)
		}
	}
	for id, worker := range document.Workers {
		if id != worker.ID {
			return fmt.Errorf("%w: worker %q", ErrInvalidDocument, id)
		}
		if err := worker.Validate(); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidDocument, err)
		}
	}
	return nil
}
