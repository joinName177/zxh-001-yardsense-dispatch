package model

import (
	"sort"
	"time"
)

type ZoneReport struct {
	Zone          string `json:"zone"`
	Open          int    `json:"open"`
	Active        int    `json:"active"`
	Resolved      int    `json:"resolved"`
	Overdue       int    `json:"overdue"`
	Unassigned    int    `json:"unassigned"`
	Critical      int    `json:"critical"`
	AverageAgeMin int    `json:"average_age_minutes"`
}

type OperationsReport struct {
	GeneratedAt   time.Time    `json:"generated_at"`
	WindowStart   time.Time    `json:"window_start"`
	WindowEnd     time.Time    `json:"window_end"`
	Zones         []ZoneReport `json:"zones"`
	TotalTickets  int          `json:"total_tickets"`
	ResolvedCount int          `json:"resolved_count"`
	OverdueCount  int          `json:"overdue_count"`
}

func BuildOperationsReport(tickets []Ticket, now time.Time, windowStart time.Time) OperationsReport {
	byZone := make(map[string]*ZoneReport)
	report := OperationsReport{GeneratedAt: now, WindowStart: windowStart, WindowEnd: now, Zones: make([]ZoneReport, 0)}
	ages := make(map[string][]int)
	for _, ticket := range tickets {
		if ticket.CreatedAt.Before(windowStart) {
			continue
		}
		report.TotalTickets++
		zone := byZone[ticket.Zone]
		if zone == nil {
			zone = &ZoneReport{Zone: ticket.Zone}
			byZone[ticket.Zone] = zone
		}
		switch ticket.Status {
		case StatusResolved:
			zone.Resolved++
			report.ResolvedCount++
		case StatusOpen:
			zone.Open++
		default:
			if !ticket.Status.Terminal() {
				zone.Active++
			}
		}
		if ticket.AssigneeID == "" && !ticket.Status.Terminal() {
			zone.Unassigned++
		}
		if ticket.Priority == PriorityCritical && !ticket.Status.Terminal() {
			zone.Critical++
		}
		if ticket.IsOverdue(now) {
			zone.Overdue++
			report.OverdueCount++
		}
		age := int(now.Sub(ticket.CreatedAt).Minutes())
		if age >= 0 {
			ages[ticket.Zone] = append(ages[ticket.Zone], age)
		}
	}
	for zoneName, summary := range byZone {
		for _, age := range ages[zoneName] {
			summary.AverageAgeMin += age
		}
		if count := len(ages[zoneName]); count > 0 {
			summary.AverageAgeMin /= count
		}
		report.Zones = append(report.Zones, *summary)
	}
	sort.Slice(report.Zones, func(i, j int) bool { return report.Zones[i].Zone < report.Zones[j].Zone })
	return report
}

type WorkloadReport struct {
	WorkerID       string `json:"worker_id"`
	WorkerName     string `json:"worker_name"`
	Capacity       int    `json:"capacity"`
	Assigned       int    `json:"assigned"`
	Available      int    `json:"available"`
	ZoneCoverage   int    `json:"zone_coverage"`
	Active         bool   `json:"active"`
	UtilizationPct int    `json:"utilization_pct"`
}

func BuildWorkloadReport(workers []Worker, tickets []Ticket) []WorkloadReport {
	loads := make(map[string]int)
	for _, ticket := range tickets {
		if ticket.AssigneeID != "" && !ticket.Status.Terminal() {
			loads[ticket.AssigneeID]++
		}
	}
	report := make([]WorkloadReport, 0, len(workers))
	for _, worker := range workers {
		assigned := loads[worker.ID]
		item := WorkloadReport{WorkerID: worker.ID, WorkerName: worker.Name, Capacity: worker.Capacity, Assigned: assigned, Available: worker.Capacity - assigned, ZoneCoverage: len(worker.Zones), Active: worker.Active}
		if item.Available < 0 {
			item.Available = 0
		}
		if worker.Capacity > 0 {
			item.UtilizationPct = assigned * 100 / worker.Capacity
		}
		report = append(report, item)
	}
	sort.Slice(report, func(i, j int) bool {
		if report[i].UtilizationPct != report[j].UtilizationPct {
			return report[i].UtilizationPct > report[j].UtilizationPct
		}
		return report[i].WorkerName < report[j].WorkerName
	})
	return report
}
