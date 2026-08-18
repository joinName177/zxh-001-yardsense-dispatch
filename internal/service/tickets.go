package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

type Clock func() time.Time

type Service struct {
	store *store.JSONStore
	now   Clock
}

func New(data *store.JSONStore, now Clock) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{store: data, now: now}
}

type CreateTicketInput struct {
	Title, Description, Zone, Operator string
	Priority                           model.Priority
	DueAt                              time.Time
	Tags                               []string
}

func (s *Service) CreateTicket(input CreateTicketInput) (model.Ticket, error) {
	if err := requireOperator(input.Operator); err != nil {
		return model.Ticket{}, err
	}
	now := s.now().UTC()
	if err := requireFutureDueAt(input.DueAt, now); err != nil {
		return model.Ticket{}, err
	}
	ticket := model.Ticket{
		ID:    generatedID("ticket", now, len(s.store.ListTickets(model.TicketFilter{}, now))+1),
		Title: input.Title, Description: input.Description, Zone: input.Zone, Priority: input.Priority,
		Status: model.StatusOpen, DueAt: input.DueAt.UTC(), CreatedAt: now, UpdatedAt: now, Revision: 1, Tags: model.NormalizeTags(input.Tags),
	}
	if err := ticket.ValidateForCreate(); err != nil {
		return model.Ticket{}, err
	}
	activity := model.Activity{TicketID: ticket.ID, Kind: model.ActivityCreated, Actor: input.Operator, Message: "ticket created", At: now}
	if err := s.store.CreateTicket(ticket, activity); err != nil {
		return model.Ticket{}, fmt.Errorf("create ticket: %w", err)
	}
	return ticket, nil
}

func (s *Service) Ticket(id string) (model.Ticket, error) {
	ticket, err := s.store.ReadTicket(id)
	if errors.Is(err, store.ErrNotFound) {
		return model.Ticket{}, ErrTicketNotFound
	}
	if err != nil {
		return model.Ticket{}, fmt.Errorf("read ticket: %w", err)
	}
	return ticket, nil
}

func (s *Service) Tickets(filter model.TicketFilter) []model.Ticket {
	return s.store.ListTickets(filter, s.now().UTC())
}

type ChangeStatusInput struct {
	TicketID, Operator string
	ExpectedRevision   int64
	Status             model.Status
}

func (s *Service) ChangeStatus(input ChangeStatusInput) (model.Ticket, error) {
	if err := requireOperator(input.Operator); err != nil {
		return model.Ticket{}, err
	}
	if !input.Status.Valid() {
		return model.Ticket{}, fmt.Errorf("invalid target status")
	}
	now := s.now().UTC()
	activity := model.Activity{TicketID: input.TicketID, Kind: model.ActivityStatus, Actor: input.Operator, Message: "status changed to " + input.Status.String(), At: now}
	ticket, err := s.store.UpdateTicket(input.TicketID, input.ExpectedRevision, func(ticket *model.Ticket) error {
		if !model.CanTransition(ticket.Status, input.Status) {
			return ErrInvalidState
		}
		ticket.Status = input.Status
		return nil
	}, activity)
	if errors.Is(err, store.ErrNotFound) {
		return model.Ticket{}, ErrTicketNotFound
	}
	if errors.Is(err, store.ErrRevisionConflict) {
		return model.Ticket{}, ErrTicketChanged
	}
	if err != nil {
		return model.Ticket{}, fmt.Errorf("change ticket status: %w", err)
	}
	return ticket, nil
}
