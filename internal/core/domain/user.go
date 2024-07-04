package domain

import (
	"database/sql"
	"strings"
)

type UserStatusType string

const (
	UserStatusUnknownStr    = "unknown"
	UserStatusActiveStr     = "ACTIVE"
	UserStatusInactiveStr   = "INACTIVE"
	UserStatusUnverifiedStr = "UNVERIFIED"
	UserStatusBannedStr     = "BANNED"
)

const (
	UserStatusUnknown    UserStatusType = UserStatusUnknownStr
	UserStatusActive     UserStatusType = UserStatusActiveStr
	UserStatusInactive   UserStatusType = UserStatusInactiveStr
	UserStatusUnverified UserStatusType = UserStatusUnverifiedStr
	UserStatusBanned     UserStatusType = UserStatusBannedStr
)

type User struct {
	Base
	Modifier

	FirstName *string
	LastName  *string
	Email     string
	Password  *string
	Avatar    *string
	Status    UserStatusType

	WelcomeMessageSent bool
	GoogleID           *string
}

func (r *User) IsActive() bool {
	return r.Status == UserStatusActive
}

func (r UserStatusType) String() string {
	var str string
	switch r {
	case UserStatusActive:
		str = UserStatusActiveStr
	case UserStatusInactive:
		str = UserStatusInactiveStr
	case UserStatusUnverified:
		str = UserStatusUnverifiedStr
	case UserStatusBanned:
		str = UserStatusBannedStr
	default:
		str = UserStatusUnknownStr
	}

	return str
}

func ToUserStatus(status string) UserStatusType {
	var userStatus UserStatusType
	switch status {
	case UserStatusActiveStr:
		userStatus = UserStatusActive
	case UserStatusInactiveStr:
		userStatus = UserStatusInactive
	case UserStatusUnverifiedStr:
		userStatus = UserStatusUnverified
	case UserStatusBannedStr:
		userStatus = UserStatusBanned
	default:
		userStatus = UserStatusUnknown
	}

	return userStatus
}

func (r *User) GetFullName() string {
	var firstName string
	if r.FirstName != nil {
		firstName = *r.FirstName
	}

	var lastName string
	if r.LastName != nil {
		lastName = *r.LastName
	}

	fullName := strings.Join([]string{firstName, lastName}, " ")
	return strings.TrimSpace(fullName)
}

func (r *User) SetGoogleID(googleID sql.NullString) *User {
	if googleID.Valid {
		r.GoogleID = &googleID.String
	}
	return r
}

func (r *User) SetFirstName(firstName sql.NullString) *User {
	if firstName.Valid {
		r.FirstName = &firstName.String
	}
	return r
}

func (r *User) SetLastName(lastName sql.NullString) *User {
	if lastName.Valid {
		r.LastName = &lastName.String
	}
	return r
}
