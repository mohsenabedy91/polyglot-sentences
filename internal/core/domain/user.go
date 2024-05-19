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
