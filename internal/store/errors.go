package store

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrAlreadyExists    = errors.New("record already exists")
	ErrRevisionConflict = errors.New("ticket revision conflict")
	ErrInvalidDocument  = errors.New("invalid persisted document")
)

type RevisionConflictError struct {
	TicketID string
	Expected int64
	Actual   int64
}

func (e *RevisionConflictError) Error() string {
	return fmt.Sprintf("ticket %s changed from revision %d to %d", e.TicketID, e.Expected, e.Actual)
}

func (e *RevisionConflictError) Is(target error) bool { return target == ErrRevisionConflict }

type MissingRecordError struct {
	Kind string
	ID   string
}

func (e *MissingRecordError) Error() string {
	return fmt.Sprintf("%s %q was not found", e.Kind, e.ID)
}

func (e *MissingRecordError) Is(target error) bool { return target == ErrNotFound }
