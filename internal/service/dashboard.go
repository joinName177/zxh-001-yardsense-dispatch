package service

import (
	"sort"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

func (s *Service) Dashboard() model.Dashboard {
	now := s.now().UTC()
	dashboard := model.NewDashboard(now)
	for _, ticket := range s.Tickets(model.TicketFilter{}) {
		dashboard.AddTicket(ticket, now)
	}
	dashboard.Recent = s.RecentActivity(12)
	sort.Slice(dashboard.Overdue, func(i, j int) bool { return dashboard.Overdue[i].DueAt.Before(dashboard.Overdue[j].DueAt) })
	return dashboard
}
