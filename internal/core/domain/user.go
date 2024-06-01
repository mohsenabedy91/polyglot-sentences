package domain

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
}

func (r User) IsActive() bool {
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
