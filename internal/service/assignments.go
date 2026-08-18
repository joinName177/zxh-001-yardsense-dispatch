package service

import (
	"errors"
	"fmt"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func (s *Service) AddWorker(worker model.Worker) error {
	if err := worker.Validate(); err != nil {
		return err
	}
	if err := s.store.UpsertWorker(worker); err != nil {
		return fmt.Errorf("save worker: %w", err)
	}
	return nil
}

func (s *Service) Workers() []model.Worker { return s.store.ListWorkers() }

func (s *Service) Assign(input model.Assignment) (model.Ticket, error) {
	if err := input.Validate(); err != nil {
		return model.Ticket{}, err
	}
	worker, err := s.store.ReadWorker(input.WorkerID)
	if errors.Is(err, store.ErrNotFound) {
		return model.Ticket{}, ErrWorkerNotFound
	}
	if err != nil {
		return model.Ticket{}, fmt.Errorf("read worker: %w", err)
	}
	if !worker.Active {
		return model.Ticket{}, ErrWorkerInactive
	}
	now := s.now().UTC()
	activity := input.Activity()
	if activity.At.IsZero() {
		activity.At = now
	}
	ticket, err := s.store.UpdateTicket(input.TicketID, input.ExpectedRev, func(ticket *model.Ticket) error {
		if ticket.Status.Terminal() {
			return ErrInvalidState
		}
		if !worker.Covers(ticket.Zone) {
			return ErrZoneMismatch
		}
		ticket.AssigneeID = worker.ID
		if ticket.Status == model.StatusOpen || ticket.Status == model.StatusBlocked {
			ticket.Status = model.StatusAssigned
		}
		return nil
	}, activity)
	if errors.Is(err, store.ErrNotFound) {
		return model.Ticket{}, ErrTicketNotFound
	}
	if errors.Is(err, store.ErrRevisionConflict) {
		return model.Ticket{}, ErrTicketChanged
	}
	if err != nil {
		return model.Ticket{}, fmt.Errorf("assign ticket: %w", err)
	}
	return ticket, nil
}
