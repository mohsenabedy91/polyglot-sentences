package domain

import (
	"github.com/google/uuid"
	"time"
)

type Modifier struct {
	CreatedBy uint
	UpdatedBy uint
}

type Base struct {
	ID   uint
	UUID uuid.UUID

	CreatedAt time.Time
	UpdatedAt time.Time
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
