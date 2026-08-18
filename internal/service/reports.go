package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func (s *Service) OperationsReport(hours int) (model.OperationsReport, error) {
	if hours < 1 || hours > 24*31 {
		return model.OperationsReport{}, fmt.Errorf("report window must be between 1 and 744 hours")
	}
	now := s.now().UTC()
	return model.BuildOperationsReport(s.Tickets(model.TicketFilter{}), now, now.Add(-time.Duration(hours)*time.Hour)), nil
}

func (s *Service) WorkloadReport() []model.WorkloadReport {
	return model.BuildWorkloadReport(s.Workers(), s.Tickets(model.TicketFilter{}))
}

func (s *Service) TicketTimeline(ticketID string) ([]model.TimelineEntry, error) {
	ticket, err := s.Ticket(ticketID)
	if err != nil {
		return nil, err
	}
	return model.BuildTimeline(ticket, s.store.RecentActivity(10000), s.store.Notes(ticketID)), nil
}

func (s *Service) ExportTicketsCSV() ([]byte, error) {
	var output bytes.Buffer
	if err := s.store.ExportTicketsCSV(&output, s.now().UTC()); err != nil {
		return nil, fmt.Errorf("export tickets: %w", err)
	}
	return output.Bytes(), nil
}

func (s *Service) ExportDocument() store.ExportDocument { return s.store.Export() }

func (s *Service) ImportDocument(document store.Document, operator string) error {
	if err := requireOperator(operator); err != nil {
		return err
	}
	if err := s.store.Import(document); err != nil {
		return fmt.Errorf("import document: %w", err)
	}
	return nil
}

func (s *Service) Backup(directory string) (store.Backup, error) { return s.store.Backup(directory) }

func (s *Service) Backups(directory string) ([]store.Backup, error) {
	return store.ListBackups(directory)
}

func (s *Service) Restore(path, operator string) error {
	if err := requireOperator(operator); err != nil {
		return err
	}
	if err := s.store.Restore(path); err != nil {
		return fmt.Errorf("restore document: %w", err)
	}
	return nil
}
