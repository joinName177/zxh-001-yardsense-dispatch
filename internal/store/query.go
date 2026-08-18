package store

import (
	"sort"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func (s *JSONStore) ListTickets(filter model.TicketFilter, now time.Time) []model.Ticket {
	document := s.Snapshot()
	tickets := make([]model.Ticket, 0, len(document.Tickets))
	for _, ticket := range document.Tickets {
		if filter.Matches(ticket, now.Unix()) {
			tickets = append(tickets, ticket.Clone())
		}
	}
	sort.Slice(tickets, func(i, j int) bool {
		if tickets[i].Priority != tickets[j].Priority {
			return tickets[i].Priority > tickets[j].Priority
		}
		return tickets[i].DueAt.Before(tickets[j].DueAt)
	})
	return tickets
}

func (s *JSONStore) ListWorkers() []model.Worker { return s.Snapshot().SortedWorkers() }

func (s *JSONStore) RecentActivity(limit int) []model.Activity {
	activities := s.Snapshot().Activities
	if limit <= 0 {
		return []model.Activity{}
	}
	if len(activities) > limit {
		activities = activities[len(activities)-limit:]
	}
	result := append([]model.Activity(nil), activities...)
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}
