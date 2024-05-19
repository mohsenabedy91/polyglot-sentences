package domain

type AccessControl struct {
	Base
	Modifier

	UserID       uint
	RoleID       *uint
	PermissionID *uint
}
