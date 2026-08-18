package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/store"
)

func (s *Service) AddNote(note model.Note) error {
	if note.Created.IsZero() {
		note.Created = s.now().UTC()
	}
	if err := note.Validate(); err != nil {
		return err
	}
	activity := model.Activity{TicketID: note.TicketID, Kind: model.ActivityNoted, Actor: note.Author, Message: "note added", At: note.Created}
	if err := s.store.AddNote(note, activity); errors.Is(err, store.ErrNotFound) {
		return ErrTicketNotFound
	} else if err != nil {
		return fmt.Errorf("save note: %w", err)
	}
	return nil
}

func (s *Service) Notes(ticketID string) ([]model.Note, error) {
	if _, err := s.Ticket(ticketID); err != nil {
		return nil, err
	}
	return s.store.Notes(ticketID), nil
}

func (s *Service) RecentActivity(limit int) []model.Activity { return s.store.RecentActivity(limit) }

func NewNote(id, ticketID, author, body string, internal bool, now time.Time) model.Note {
	return model.Note{ID: id, TicketID: ticketID, Author: author, Body: body, Internal: internal, Created: now.UTC()}
}
