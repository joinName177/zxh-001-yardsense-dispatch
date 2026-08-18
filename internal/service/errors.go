package service

import "errors"

var (
	ErrTicketNotFound = errors.New("ticket not found")
	ErrWorkerNotFound = errors.New("worker not found")
	ErrTicketChanged  = errors.New("ticket changed by another dispatcher")
	ErrWorkerInactive = errors.New("worker is inactive")
	ErrZoneMismatch   = errors.New("worker does not cover ticket zone")
	ErrInvalidState   = errors.New("invalid ticket state transition")
)
