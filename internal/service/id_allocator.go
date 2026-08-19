package service

import (
	"fmt"
	"sync/atomic"
	"time"
)

type ticketIDAllocator struct {
	counter atomic.Int64
}

func newTicketIDAllocator() *ticketIDAllocator { return &ticketIDAllocator{} }

func (a *ticketIDAllocator) Next(now time.Time) string {
	seq := a.counter.Add(1)
	return fmt.Sprintf("ticket-%d-%03d", now.UnixNano(), seq)
}
