package model

import "fmt"

// Status describes the operational stage of a yard ticket.
type Status string

const (
	StatusOpen       Status = "open"
	StatusAssigned   Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusBlocked    Status = "blocked"
	StatusResolved   Status = "resolved"
	StatusCancelled  Status = "cancelled"
)

var validStatuses = map[Status]struct{}{
	StatusOpen: {}, StatusAssigned: {}, StatusInProgress: {},
	StatusBlocked: {}, StatusResolved: {}, StatusCancelled: {},
}

func (s Status) Valid() bool {
	_, ok := validStatuses[s]
	return ok
}

func (s Status) Terminal() bool {
	return s == StatusResolved || s == StatusCancelled
}

func CanTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusOpen:
		return to == StatusAssigned || to == StatusCancelled
	case StatusAssigned:
		return to == StatusInProgress || to == StatusBlocked || to == StatusCancelled
	case StatusInProgress:
		return to == StatusBlocked || to == StatusResolved || to == StatusCancelled
	case StatusBlocked:
		return to == StatusAssigned || to == StatusInProgress || to == StatusCancelled
	default:
		return false
	}
}

func (s Status) String() string { return string(s) }

func ParseStatus(value string) (Status, error) {
	status := Status(value)
	if !status.Valid() {
		return "", fmt.Errorf("unknown ticket status %q", value)
	}
	return status, nil
}
