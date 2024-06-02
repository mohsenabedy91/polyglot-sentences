package domain

import "database/sql"

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
	UserStatusInActive   UserStatusType = UserStatusInactiveStr
	UserStatusUnVerified UserStatusType = UserStatusUnverifiedStr
	UserStatusBanned     UserStatusType = UserStatusBannedStr
)

type User struct {
	Base
	Modifier

	FirstName string
	LastName  string
	Email     string
	Password  string
	Avatar    string
	Status    UserStatusType

	WelcomeMessageSent bool
	GoogleID           string
}

func (r *User) IsActive() bool {
	return r.Status == UserStatusActive
}

func (r UserStatusType) String() string {
	var str string
	switch r {
	case UserStatusActive:
		str = UserStatusActiveStr
	case UserStatusInActive:
		str = UserStatusInactiveStr
	case UserStatusUnVerified:
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
		userStatus = UserStatusInActive
	case UserStatusUnverifiedStr:
		userStatus = UserStatusUnVerified
	case UserStatusBannedStr:
		userStatus = UserStatusBanned
	default:
		userStatus = UserStatusUnknown
	}

	return userStatus
}

func (r *User) GetFullName() string {
	return r.FirstName + " " + r.LastName
}

func (r *User) SetGoogleID(googleID sql.NullString) *User {
	if googleID.Valid {
		r.GoogleID = googleID.String
	}
	return r
}

func (r *User) SetFirstName(firstName sql.NullString) *User {
	if firstName.Valid {
		r.FirstName = firstName.String
	}
	return r
}

func (r *User) SetLastName(lastName sql.NullString) *User {
	if lastName.Valid {
		r.LastName = lastName.String
	}
	return r
}
