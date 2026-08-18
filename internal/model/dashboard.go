package model

import "time"

type Dashboard struct {
	GeneratedAt time.Time      `json:"generated_at"`
	ByStatus    map[Status]int `json:"by_status"`
	ByZone      map[string]int `json:"by_zone"`
	Overdue     []Ticket       `json:"overdue"`
	WorkerLoad  map[string]int `json:"worker_load"`
	Recent      []Activity     `json:"recent_activity"`
}

func NewDashboard(now time.Time) Dashboard {
	return Dashboard{
		GeneratedAt: now, ByStatus: make(map[Status]int), ByZone: make(map[string]int),
		WorkerLoad: make(map[string]int), Overdue: make([]Ticket, 0), Recent: make([]Activity, 0),
	}
}

func (d *Dashboard) AddTicket(ticket Ticket, now time.Time) {
	d.ByStatus[ticket.Status]++
	d.ByZone[ticket.Zone]++
	if ticket.AssigneeID != "" && !ticket.Status.Terminal() {
		d.WorkerLoad[ticket.AssigneeID]++
	}
	if ticket.IsOverdue(now) {
		d.Overdue = append(d.Overdue, ticket.Clone())
	}
}
