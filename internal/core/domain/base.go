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
	StatusActive      StatusType = "ACTIVE"
	StatusDisabled    StatusType = "DISABLED"
	StatusUnpublished StatusType = "UNPUBLISHED"
	StatusDraft       StatusType = "DRAFT"
)

type UserStatusType string

const (
	UserStatusActive     UserStatusType = "ACTIVE"
	UserStatusInActive   UserStatusType = "INACTIVE"
	UserStatusUnVerified UserStatusType = "UNVERIFIED"
	UserStatusBanned     UserStatusType = "BANNED"
)
