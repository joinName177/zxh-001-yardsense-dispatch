package store

import (
	"sort"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

type Document struct {
	Version    int                     `json:"version"`
	Tickets    map[string]model.Ticket `json:"tickets"`
	Workers    map[string]model.Worker `json:"workers"`
	Notes      map[string][]model.Note `json:"notes"`
	Activities []model.Activity        `json:"activities"`
	UpdatedAt  time.Time               `json:"updated_at"`
}

func NewDocument() Document {
	return Document{
		Version: 1, Tickets: make(map[string]model.Ticket), Workers: make(map[string]model.Worker),
		Notes: make(map[string][]model.Note), Activities: make([]model.Activity, 0),
	}
}

func (d Document) Clone() Document {
	copy := NewDocument()
	copy.Version, copy.UpdatedAt = d.Version, d.UpdatedAt
	for id, ticket := range d.Tickets {
		copy.Tickets[id] = ticket.Clone()
	}
	for id, worker := range d.Workers {
		copy.Workers[id] = worker.Clone()
	}
	for id, notes := range d.Notes {
		copy.Notes[id] = append([]model.Note(nil), notes...)
	}
	copy.Activities = append([]model.Activity(nil), d.Activities...)
	return copy
}

func (d Document) SortedTickets() []model.Ticket {
	tickets := make([]model.Ticket, 0, len(d.Tickets))
	for _, ticket := range d.Tickets {
		tickets = append(tickets, ticket.Clone())
	}
	sort.Slice(tickets, func(i, j int) bool { return tickets[i].CreatedAt.After(tickets[j].CreatedAt) })
	return tickets
}

func (d Document) SortedWorkers() []model.Worker {
	workers := make([]model.Worker, 0, len(d.Workers))
	for _, worker := range d.Workers {
		workers = append(workers, worker.Clone())
	}
	model.SortWorkers(workers)
	return workers
}
