package service

import (
	"fmt"
	"sync"
	"time"
)

type ticketIDAllocator struct {
	mu       sync.Mutex
	lastNano int64
	sequence int
}

func newTicketIDAllocator() *ticketIDAllocator { return &ticketIDAllocator{} }

func (a *ticketIDAllocator) Next(now time.Time) string {
	a.mu.Lock()
	defer a.mu.Unlock()
	nano := now.UnixNano()
	if nano != a.lastNano {
		a.lastNano, a.sequence = nano, 0
	}
	a.sequence++
	return fmt.Sprintf("ticket-%d-%03d", nano, a.sequence)
}
