package domain

import (
	"github.com/google/uuid"
	"time"
)

type Modifier struct {
	CreatedBy *uint64
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

type ChallengeType string

const UnknownStr = "UNKNOWN"

const (
	ChallengeTOTPStr          = "TOTP"
	ChallengeRecoveryCodesStr = "RECOVERY_CODES"
)

const (
	ChallengeTOTP          ChallengeType = ChallengeTOTPStr
	ChallengeRecoveryCodes ChallengeType = ChallengeRecoveryCodesStr
)

type ChallengeStatusType string

const (
	ChallengeEnableStr   = "ENABLE"
	ChallengeDisabledStr = "DISABLED"
)

const (
	ChallengeEnable   ChallengeStatusType = "ENABLE"
	ChallengeDisabled ChallengeStatusType = "DISABLED"
)

type AuthChallenge struct {
	Type   ChallengeType
	Status ChallengeStatusType
}

func (r ChallengeType) String() string {
	var str string
	switch r {
	case ChallengeTOTP:
		str = ChallengeTOTPStr
	case ChallengeRecoveryCodes:
		str = ChallengeRecoveryCodesStr
	default:
		str = UnknownStr
	}

	return str
}

func (r ChallengeStatusType) String() string {
	var str string
	switch r {
	case ChallengeEnable:
		str = ChallengeEnableStr
	case ChallengeDisabled:
		str = ChallengeDisabledStr
	default:
		str = UnknownStr
	}

	return str
}
