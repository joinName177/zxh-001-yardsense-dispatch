package model

import (
	"fmt"
	"time"
)

type ActivityKind string

const (
	ActivityCreated  ActivityKind = "created"
	ActivityAssigned ActivityKind = "assigned"
	ActivityStatus   ActivityKind = "status_changed"
	ActivityNoted    ActivityKind = "noted"
)

type Activity struct {
	TicketID string       `json:"ticket_id"`
	Kind     ActivityKind `json:"kind"`
	Actor    string       `json:"actor"`
	Message  string       `json:"message"`
	At       time.Time    `json:"at"`
}

func (a Activity) Validate() error {
	if a.TicketID == "" || a.Actor == "" || a.Message == "" || a.At.IsZero() {
		return fmt.Errorf("activity requires ticket, actor, message, and time")
	}
	switch a.Kind {
	case ActivityCreated, ActivityAssigned, ActivityStatus, ActivityNoted:
		return nil
	default:
		return fmt.Errorf("unknown activity kind %q", a.Kind)
	}
}

func (a Activity) Clone() Activity { return a }
