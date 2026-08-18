package service

import (
	"fmt"
	"strings"
	"time"
)

func requireOperator(operator string) error {
	if strings.TrimSpace(operator) == "" {
		return fmt.Errorf("operator is required")
	}
	return nil
}

func requireFutureDueAt(dueAt, now time.Time) error {
	if dueAt.Before(now.Add(-5 * time.Minute)) {
		return fmt.Errorf("due time cannot be more than five minutes in the past")
	}
	return nil
}

func generatedID(prefix string, now time.Time, sequence int) string {
	return fmt.Sprintf("%s-%d-%03d", prefix, now.UnixNano(), sequence)
}
