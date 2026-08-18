package model

import (
	"sort"
	"strings"
)

type AssignmentRecommendation struct {
	WorkerID    string   `json:"worker_id"`
	WorkerName  string   `json:"worker_name"`
	Score       int      `json:"score"`
	Reasons     []string `json:"reasons"`
	CurrentLoad int      `json:"current_load"`
	Remaining   int      `json:"remaining_capacity"`
	Eligible    bool     `json:"eligible"`
}

func RecommendWorkers(ticket Ticket, workers []Worker, activeLoad map[string]int) []AssignmentRecommendation {
	recommendations := make([]AssignmentRecommendation, 0, len(workers))
	for _, worker := range workers {
		load := activeLoad[worker.ID]
		item := AssignmentRecommendation{WorkerID: worker.ID, WorkerName: worker.Name, CurrentLoad: load, Remaining: worker.Capacity - load, Reasons: make([]string, 0, 4)}
		if !worker.Active {
			item.Reasons = append(item.Reasons, "worker is inactive")
			recommendations = append(recommendations, item)
			continue
		}
		if !worker.Covers(ticket.Zone) {
			item.Reasons = append(item.Reasons, "does not cover ticket zone")
			recommendations = append(recommendations, item)
			continue
		}
		if load >= worker.Capacity {
			item.Reasons = append(item.Reasons, "at capacity")
			recommendations = append(recommendations, item)
			continue
		}
		item.Eligible = true
		item.Score = 50 + item.Remaining*8
		item.Reasons = append(item.Reasons, "covers "+ticket.Zone, "has "+itoa(item.Remaining)+" open capacity")
		for _, tag := range ticket.Tags {
			if worker.HasSkill(tag) {
				item.Score += 12
				item.Reasons = append(item.Reasons, "matches skill "+tag)
			}
		}
		if ticket.Priority == PriorityCritical && worker.HasSkill("safety") {
			item.Score += 15
			item.Reasons = append(item.Reasons, "safety qualified for critical work")
		}
		recommendations = append(recommendations, item)
	}
	sort.SliceStable(recommendations, func(i, j int) bool {
		if recommendations[i].Eligible != recommendations[j].Eligible {
			return recommendations[i].Eligible
		}
		if recommendations[i].Score != recommendations[j].Score {
			return recommendations[i].Score > recommendations[j].Score
		}
		return strings.ToLower(recommendations[i].WorkerName) < strings.ToLower(recommendations[j].WorkerName)
	})
	return recommendations
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := [20]byte{}
	position := len(digits)
	for value > 0 {
		position--
		digits[position] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[position:])
}
