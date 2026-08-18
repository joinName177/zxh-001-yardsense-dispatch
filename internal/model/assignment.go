package model

import (
	"fmt"
	"time"
)

type Assignment struct {
	TicketID    string    `json:"ticket_id"`
	WorkerID    string    `json:"worker_id"`
	AssignedBy  string    `json:"assigned_by"`
	AssignedAt  time.Time `json:"assigned_at"`
	ExpectedRev int64     `json:"expected_revision"`
}

func (a Assignment) Validate() error {
	if a.TicketID == "" {
		return fmt.Errorf("assignment ticket id is required")
	}
	if a.WorkerID == "" {
		return fmt.Errorf("assignment worker id is required")
	}
	if a.AssignedBy == "" {
		return fmt.Errorf("assignment operator is required")
	}
	if a.ExpectedRev < 1 {
		return fmt.Errorf("assignment expected revision must be positive")
	}
	return nil
}

func (a Assignment) Activity() Activity {
	return Activity{
		TicketID: a.TicketID, Kind: ActivityAssigned, Actor: a.AssignedBy,
		Message: fmt.Sprintf("assigned to %s", a.WorkerID), At: a.AssignedAt,
	}
}
