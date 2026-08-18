package model

import (
	"sort"
	"time"
)

type TimelineEntry struct {
	At       time.Time `json:"at"`
	TicketID string    `json:"ticket_id"`
	Kind     string    `json:"kind"`
	Actor    string    `json:"actor"`
	Detail   string    `json:"detail"`
}

func BuildTimeline(ticket Ticket, activities []Activity, notes []Note) []TimelineEntry {
	timeline := make([]TimelineEntry, 0, len(activities)+len(notes))
	for _, activity := range activities {
		if activity.TicketID != ticket.ID {
			continue
		}
		timeline = append(timeline, TimelineEntry{At: activity.At, TicketID: ticket.ID, Kind: string(activity.Kind), Actor: activity.Actor, Detail: activity.Message})
	}
	for _, note := range notes {
		timeline = append(timeline, TimelineEntry{At: note.Created, TicketID: ticket.ID, Kind: "note", Actor: note.Author, Detail: note.Preview(200)})
	}
	sort.Slice(timeline, func(i, j int) bool { return timeline[i].At.Before(timeline[j].At) })
	return timeline
}

type Handoff struct {
	TicketID     string    `json:"ticket_id"`
	FromOperator string    `json:"from_operator"`
	ToOperator   string    `json:"to_operator"`
	Summary      string    `json:"summary"`
	Acknowledged bool      `json:"acknowledged"`
	CreatedAt    time.Time `json:"created_at"`
}

func (h Handoff) Ready() bool {
	return h.TicketID != "" && h.FromOperator != "" && h.ToOperator != "" && h.Summary != ""
}
