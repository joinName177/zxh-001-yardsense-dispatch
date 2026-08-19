package service

import (
	"fmt"
	"time"
)

type ticketIDAllocator struct{}

func newTicketIDAllocator() *ticketIDAllocator { return &ticketIDAllocator{} }

func (a *ticketIDAllocator) Next(now time.Time) string {
	return fmt.Sprintf("ticket-%d-%03d", now.UnixNano(), 1)
}
