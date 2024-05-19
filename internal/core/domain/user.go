package domain

type User struct {
	Base
	Modifier

	FirstName *string
	LastName  *string
	Email     *string
	Password  *string
	Avatar    *string
	Status    UserStatusType
}

func (r User) IsActive() bool {
	return r.Status == UserStatusActive
}

func (r UserStatusType) String() *string {
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
	}

	return &str
}
