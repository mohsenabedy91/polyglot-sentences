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

type UserStatusType string

const (
	UserStatusActiveStr     = "ACTIVE"
	UserStatusInactiveStr   = "INACTIVE"
	UserStatusUnverifiedStr = "UNVERIFIED"
	UserStatusBannedStr     = "BANNED"
)

const (
	UserStatusActive     UserStatusType = UserStatusActiveStr
	UserStatusInActive   UserStatusType = UserStatusInactiveStr
	UserStatusUnVerified UserStatusType = UserStatusUnverifiedStr
	UserStatusBanned     UserStatusType = UserStatusBannedStr
)
