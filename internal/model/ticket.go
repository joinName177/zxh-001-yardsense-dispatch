package model

import (
	"fmt"
	"strings"
	"time"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

func (p Priority) Valid() bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Zone        string    `json:"zone"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	AssigneeID  string    `json:"assignee_id,omitempty"`
	DueAt       time.Time `json:"due_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Revision    int64     `json:"revision"`
	Tags        []string  `json:"tags"`
}

func (t Ticket) ValidateForCreate() error {
	if len(strings.TrimSpace(t.Title)) < 4 {
		return fmt.Errorf("ticket title must be at least 4 characters")
	}
	if len(strings.TrimSpace(t.Description)) < 10 {
		return fmt.Errorf("ticket description must be at least 10 characters")
	}
	if len(strings.TrimSpace(t.Zone)) == 0 {
		return fmt.Errorf("ticket zone is required")
	}
	if !t.Priority.Valid() {
		return fmt.Errorf("ticket priority %q is invalid", t.Priority)
	}
	if t.DueAt.IsZero() {
		return fmt.Errorf("ticket due time is required")
	}
	return nil
}

func (t Ticket) Clone() Ticket {
	copy := t
	copy.Tags = append([]string(nil), t.Tags...)
	return copy
}

func (t Ticket) IsOverdue(now time.Time) bool {
	return !t.Status.Terminal() && t.DueAt.Before(now)
}

func (t Ticket) Summary() string {
	assignee := t.AssigneeID
	if assignee == "" {
		assignee = "unassigned"
	}
	return fmt.Sprintf("%s [%s] in %s for %s", t.Title, t.Status, t.Zone, assignee)
}

func NormalizeTags(tags []string) []string {
	unique := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		value := strings.ToLower(strings.TrimSpace(tag))
		if value == "" {
			continue
		}
		if _, ok := unique[value]; ok {
			continue
		}
		unique[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
