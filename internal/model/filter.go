package model

import "strings"

type TicketFilter struct {
	Status     Status
	Zone       string
	AssigneeID string
	Tag        string
	Overdue    bool
}

func (f TicketFilter) Matches(ticket Ticket, nowUnix int64) bool {
	if f.Status != "" && ticket.Status != f.Status {
		return false
	}
	if f.Zone != "" && !strings.EqualFold(ticket.Zone, f.Zone) {
		return false
	}
	if f.AssigneeID != "" && ticket.AssigneeID != f.AssigneeID {
		return false
	}
	if f.Tag != "" {
		matched := false
		for _, tag := range ticket.Tags {
			if strings.EqualFold(tag, f.Tag) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return !f.Overdue || (!ticket.Status.Terminal() && ticket.DueAt.Unix() < nowUnix)
}
