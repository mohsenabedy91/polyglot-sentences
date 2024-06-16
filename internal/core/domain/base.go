package domain

import (
	"github.com/google/uuid"
	"time"
)

type Modifier struct {
	CreatedBy uint64
	UpdatedBy uint64
	DeleteBy  uint64
}

type Base struct {
	ID   uint64
	UUID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type StatusType string

const (
	StatusActiveStr      = "ACTIVE"
	StatusDisabledStr    = "DISABLED"
	StatusUnpublishedStr = "UNPUBLISHED"
	StatusDraftStr       = "DRAFT"
)

const (
	StatusActive      StatusType = StatusActiveStr
	StatusDisabled    StatusType = StatusDisabledStr
	StatusUnpublished StatusType = StatusUnpublishedStr
	StatusDraft       StatusType = StatusDraftStr
)

func (r StatusType) String() *string {
	var str string
	switch r {
	case StatusActive:
		str = StatusActiveStr
	case StatusDisabled:
		str = StatusDisabledStr
	case StatusUnpublished:
		str = StatusUnpublishedStr
	case StatusDraft:
		str = StatusDraftStr
	}

	return &str
}
