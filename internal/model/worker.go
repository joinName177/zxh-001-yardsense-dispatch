package model

import (
	"fmt"
	"sort"
	"strings"
)

type Worker struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Zones    []string `json:"zones"`
	Skills   []string `json:"skills"`
	Active   bool     `json:"active"`
	Capacity int      `json:"capacity"`
}

func (w Worker) Validate() error {
	if strings.TrimSpace(w.ID) == "" || strings.TrimSpace(w.Name) == "" {
		return fmt.Errorf("worker id and name are required")
	}
	if w.Capacity < 1 || w.Capacity > 20 {
		return fmt.Errorf("worker capacity must be between 1 and 20")
	}
	if len(w.Zones) == 0 {
		return fmt.Errorf("worker must cover at least one zone")
	}
	return nil
}

func (w Worker) Covers(zone string) bool {
	for _, candidate := range w.Zones {
		if strings.EqualFold(candidate, zone) {
			return true
		}
	}
	return false
}

func (w Worker) HasSkill(skill string) bool {
	for _, candidate := range w.Skills {
		if strings.EqualFold(candidate, skill) {
			return true
		}
	}
	return false
}

func (w Worker) Clone() Worker {
	copy := w
	copy.Zones = append([]string(nil), w.Zones...)
	copy.Skills = append([]string(nil), w.Skills...)
	return copy
}

func SortWorkers(workers []Worker) {
	sort.Slice(workers, func(i, j int) bool {
		return strings.ToLower(workers[i].Name) < strings.ToLower(workers[j].Name)
	})
}
