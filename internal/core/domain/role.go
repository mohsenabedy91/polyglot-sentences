package domain

type RoleKeyType string

const (
	RoleKeySuperAdmin RoleKeyType = "SUPER_ADMIN"
	RoleKeyAdmin      RoleKeyType = "ADMIN"
	RoleKeyManager    RoleKeyType = "MANAGER"
	RoleKeyAccountant RoleKeyType = "ACCOUNTANT"
	RoleKeySupplier   RoleKeyType = "SUPPLIER"
	RoleKeySales      RoleKeyType = "SALES"
	RoleKeyStaff      RoleKeyType = "STAFF"
	RoleKeyUser       RoleKeyType = "USER"
)

type Role struct {
	Base
	Modifier

	Title       string
	Key         RoleKeyType
	Description string
}
