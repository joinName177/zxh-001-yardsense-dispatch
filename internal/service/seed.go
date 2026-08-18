package service

import (
	"fmt"
	"time"

	"github.com/joinName177/zxh-001-yardsense-dispatch/internal/model"
)

type SeedResult struct {
	Workers int `json:"workers"`
	Tickets int `json:"tickets"`
}

func (s *Service) SeedDemo(operator string) (SeedResult, error) {
	if err := requireOperator(operator); err != nil {
		return SeedResult{}, err
	}
	workers := []model.Worker{
		{ID: "worker-maya", Name: "Maya Chen", Zones: []string{"north-yard", "gate-a"}, Skills: []string{"safety", "equipment"}, Active: true, Capacity: 4},
		{ID: "worker-alex", Name: "Alex Park", Zones: []string{"south-yard", "gate-b"}, Skills: []string{"loading", "equipment"}, Active: true, Capacity: 3},
		{ID: "worker-riley", Name: "Riley Sun", Zones: []string{"north-yard", "south-yard"}, Skills: []string{"safety", "loading"}, Active: true, Capacity: 5},
	}
	result := SeedResult{}
	for _, worker := range workers {
		if err := s.AddWorker(worker); err != nil {
			return result, fmt.Errorf("seed worker %s: %w", worker.ID, err)
		}
		result.Workers++
	}
	if len(s.Tickets(model.TicketFilter{})) > 0 {
		return result, nil
	}
	now := s.now().UTC()
	inputs := []CreateTicketInput{
		{Title: "Inspect damaged pallet rack", Description: "A forklift operator reported movement in bay N-14 during unloading.", Zone: "north-yard", Priority: model.PriorityHigh, DueAt: now.Add(2 * time.Hour), Operator: operator, Tags: []string{"safety", "rack"}},
		{Title: "Clear gate B sensor fault", Description: "The inbound sensor is intermittently failing to register trailer entries.", Zone: "gate-b", Priority: model.PriorityNormal, DueAt: now.Add(4 * time.Hour), Operator: operator, Tags: []string{"gate", "equipment"}},
		{Title: "Review spill response stock", Description: "Night shift requested a replenishment check for spill response supplies.", Zone: "south-yard", Priority: model.PriorityLow, DueAt: now.Add(8 * time.Hour), Operator: operator, Tags: []string{"safety", "inventory"}},
	}
	for _, input := range inputs {
		if _, err := s.CreateTicket(input); err != nil {
			return result, fmt.Errorf("seed ticket: %w", err)
		}
		result.Tickets++
	}
	return result, nil
}
