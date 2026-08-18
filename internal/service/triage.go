package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func (s *Service) RecommendAssignment(ticketID string) ([]model.AssignmentRecommendation, error) {
	ticket, err := s.Ticket(ticketID)
	if err != nil {
		return nil, err
	}
	if ticket.Status.Terminal() {
		return nil, fmt.Errorf("cannot recommend a worker for terminal ticket")
	}
	loads := make(map[string]int)
	for _, candidate := range s.Tickets(model.TicketFilter{}) {
		if candidate.AssigneeID != "" && !candidate.Status.Terminal() {
			loads[candidate.AssigneeID]++
		}
	}
	return model.RecommendWorkers(ticket, s.Workers(), loads), nil
}

type TriageQueueItem struct {
	Ticket         model.Ticket `json:"ticket"`
	UrgencyScore   int          `json:"urgency_score"`
	RecommendedFor string       `json:"recommended_for,omitempty"`
	Reason         string       `json:"reason"`
}

func (s *Service) TriageQueue() []TriageQueueItem {
	now := s.now().UTC()
	queue := make([]TriageQueueItem, 0)
	for _, ticket := range s.Tickets(model.TicketFilter{}) {
		if ticket.Status.Terminal() {
			continue
		}
		item := TriageQueueItem{Ticket: ticket, UrgencyScore: urgency(ticket, now), Reason: explainUrgency(ticket, now)}
		if ticket.AssigneeID == "" {
			if recommendations, err := s.RecommendAssignment(ticket.ID); err == nil && len(recommendations) > 0 && recommendations[0].Eligible {
				item.RecommendedFor = recommendations[0].WorkerID
			}
		}
		queue = append(queue, item)
	}
	sort.Slice(queue, func(i, j int) bool {
		if queue[i].UrgencyScore != queue[j].UrgencyScore {
			return queue[i].UrgencyScore > queue[j].UrgencyScore
		}
		return queue[i].Ticket.DueAt.Before(queue[j].Ticket.DueAt)
	})
	return queue
}

func urgency(ticket model.Ticket, now time.Time) int {
	score := 0
	switch ticket.Priority {
	case model.PriorityCritical:
		score += 100
	case model.PriorityHigh:
		score += 70
	case model.PriorityNormal:
		score += 40
	default:
		score += 15
	}
	if ticket.AssigneeID == "" {
		score += 25
	}
	if ticket.Status == model.StatusBlocked {
		score += 20
	}
	untilDue := ticket.DueAt.Sub(now)
	switch {
	case untilDue < 0:
		score += 50
	case untilDue < time.Hour:
		score += 35
	case untilDue < 4*time.Hour:
		score += 15
	}
	return score
}

func explainUrgency(ticket model.Ticket, now time.Time) string {
	reasons := []string{string(ticket.Priority) + " priority"}
	if ticket.AssigneeID == "" {
		reasons = append(reasons, "unassigned")
	}
	if ticket.Status == model.StatusBlocked {
		reasons = append(reasons, "blocked")
	}
	if ticket.DueAt.Before(now) {
		reasons = append(reasons, "overdue")
	} else if ticket.DueAt.Before(now.Add(time.Hour)) {
		reasons = append(reasons, "due within one hour")
	}
	return strings.Join(reasons, ", ")
}

type TriageSummary struct {
	Total      int `json:"total"`
	Critical   int `json:"critical"`
	Unassigned int `json:"unassigned"`
	Overdue    int `json:"overdue"`
	Blocked    int `json:"blocked"`
}

func (s *Service) TriageSummary() TriageSummary {
	now := s.now().UTC()
	summary := TriageSummary{}
	for _, item := range s.TriageQueue() {
		summary.Total++
		if item.Ticket.Priority == model.PriorityCritical {
			summary.Critical++
		}
		if item.Ticket.AssigneeID == "" {
			summary.Unassigned++
		}
		if item.Ticket.IsOverdue(now) {
			summary.Overdue++
		}
		if item.Ticket.Status == model.StatusBlocked {
			summary.Blocked++
		}
	}
	return summary
}

func (s *Service) EligibleWorkers(ticketID string) ([]model.Worker, error) {
	recommendations, err := s.RecommendAssignment(ticketID)
	if err != nil {
		return nil, err
	}
	workers := make([]model.Worker, 0, len(recommendations))
	byID := make(map[string]model.Worker)
	for _, worker := range s.Workers() {
		byID[worker.ID] = worker
	}
	for _, recommendation := range recommendations {
		if recommendation.Eligible {
			workers = append(workers, byID[recommendation.WorkerID])
		}
	}
	return workers, nil
}
