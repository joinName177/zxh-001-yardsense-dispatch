package model

import (
	"fmt"
	"strings"
	"time"
)

type Note struct {
	ID       string    `json:"id"`
	TicketID string    `json:"ticket_id"`
	Author   string    `json:"author"`
	Body     string    `json:"body"`
	Created  time.Time `json:"created"`
	Internal bool      `json:"internal"`
}

func (n Note) Validate() error {
	if n.TicketID == "" || n.Author == "" {
		return fmt.Errorf("note ticket id and author are required")
	}
	if len(strings.TrimSpace(n.Body)) < 3 {
		return fmt.Errorf("note body must contain at least 3 characters")
	}
	if len(n.Body) > 4000 {
		return fmt.Errorf("note body exceeds 4000 characters")
	}
	return nil
}

func (n Note) Preview(limit int) string {
	value := strings.TrimSpace(n.Body)
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "..."
}
