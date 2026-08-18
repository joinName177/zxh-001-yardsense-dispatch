package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

type JSONStore struct {
	path string
	mu   sync.RWMutex
	doc  Document
}

func Open(path string) (*JSONStore, error) {
	store := &JSONStore{path: path, doc: NewDocument()}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *JSONStore) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read data file: %w", err)
	}
	if err := json.Unmarshal(data, &s.doc); err != nil {
		return fmt.Errorf("decode data file: %w", err)
	}
	if s.doc.Version != 1 || s.doc.Tickets == nil || s.doc.Workers == nil || s.doc.Notes == nil {
		return ErrInvalidDocument
	}
	return nil
}

func (s *JSONStore) persistLocked() error {
	s.doc.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(s.doc, "", "  ")
	if err != nil {
		return fmt.Errorf("encode document: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	temporary := s.path + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return fmt.Errorf("write temporary document: %w", err)
	}
	if err := os.Rename(temporary, s.path); err != nil {
		return fmt.Errorf("replace document: %w", err)
	}
	return nil
}

func (s *JSONStore) Snapshot() Document {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.doc.Clone()
}

func (s *JSONStore) CreateTicket(ticket model.Ticket, activity model.Activity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.doc.Tickets[ticket.ID]; exists {
		return ErrAlreadyExists
	}
	s.doc.Tickets[ticket.ID] = ticket.Clone()
	s.doc.Activities = append(s.doc.Activities, activity.Clone())
	return s.persistLocked()
}

func (s *JSONStore) ReadTicket(id string) (model.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ticket, exists := s.doc.Tickets[id]
	if !exists {
		return model.Ticket{}, &MissingRecordError{Kind: "ticket", ID: id}
	}
	return ticket.Clone(), nil
}

func (s *JSONStore) UpdateTicket(id string, expectedRevision int64, mutate func(*model.Ticket) error, activity model.Activity) (model.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ticket, exists := s.doc.Tickets[id]
	if !exists {
		return model.Ticket{}, &MissingRecordError{Kind: "ticket", ID: id}
	}
	if ticket.Revision != expectedRevision {
		conflict := &RevisionConflictError{TicketID: id, Expected: expectedRevision, Actual: ticket.Revision}
		return model.Ticket{}, fmt.Errorf("update ticket rejected: %w", conflict)
	}
	updated := ticket.Clone()
	if err := mutate(&updated); err != nil {
		return model.Ticket{}, err
	}
	updated.Revision++
	updated.UpdatedAt = time.Now().UTC()
	s.doc.Tickets[id] = updated
	s.doc.Activities = append(s.doc.Activities, activity.Clone())
	if err := s.persistLocked(); err != nil {
		s.doc.Tickets[id] = ticket
		s.doc.Activities = s.doc.Activities[:len(s.doc.Activities)-1]
		return model.Ticket{}, err
	}
	return updated.Clone(), nil
}

func (s *JSONStore) UpsertWorker(worker model.Worker) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.doc.Workers[worker.ID] = worker.Clone()
	return s.persistLocked()
}

func (s *JSONStore) ReadWorker(id string) (model.Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	worker, exists := s.doc.Workers[id]
	if !exists {
		return model.Worker{}, &MissingRecordError{Kind: "worker", ID: id}
	}
	return worker.Clone(), nil
}

func (s *JSONStore) AddNote(note model.Note, activity model.Activity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.doc.Tickets[note.TicketID]; !exists {
		return &MissingRecordError{Kind: "ticket", ID: note.TicketID}
	}
	s.doc.Notes[note.TicketID] = append(s.doc.Notes[note.TicketID], note)
	s.doc.Activities = append(s.doc.Activities, activity)
	return s.persistLocked()
}

func (s *JSONStore) Notes(ticketID string) []model.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]model.Note(nil), s.doc.Notes[ticketID]...)
}
