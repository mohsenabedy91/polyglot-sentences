package domain

type PermissionKeyType string

const (
	PermissionKeyCreateUser           PermissionKeyType = "CREATE_USER"
	PermissionKeyReadUser             PermissionKeyType = "READ_USER"
	PermissionKeyUpdateUser           PermissionKeyType = "UPDATE_USER"
	PermissionKeyDeleteUser           PermissionKeyType = "DELETE_USER"
	PermissionKeyCreateRole           PermissionKeyType = "CREATE_ROLE"
	PermissionKeyReadRole             PermissionKeyType = "READ_ROLE"
	PermissionKeyUpdateRole           PermissionKeyType = "UPDATE_ROLE"
	PermissionKeyDeleteRole           PermissionKeyType = "DELETE_ROLE"
	PermissionKeyCreatePermission     PermissionKeyType = "CREATE_PERMISSION"
	PermissionKeyReadPermission       PermissionKeyType = "READ_PERMISSION"
	PermissionKeyUpdatePermission     PermissionKeyType = "UPDATE_PERMISSION"
	PermissionKeyDeletePermission     PermissionKeyType = "DELETE_PERMISSION"
	PermissionKeyCreateUserRole       PermissionKeyType = "CREATE_USER_ROLE"
	PermissionKeyReadUserRole         PermissionKeyType = "READ_USER_ROLE"
	PermissionKeyUpdateUserRole       PermissionKeyType = "UPDATE_USER_ROLE"
	PermissionKeyDeleteUserRole       PermissionKeyType = "DELETE_USER_ROLE"
	PermissionKeyCreateRolePermission PermissionKeyType = "CREATE_ROLE_PERMISSION"
	PermissionKeyReadRolePermission   PermissionKeyType = "READ_ROLE_PERMISSION"
	PermissionKeyUpdateRolePermission PermissionKeyType = "UPDATE_ROLE_PERMISSION"
	PermissionKeyDeleteRolePermission PermissionKeyType = "DELETE_ROLE_PERMISSION"
)

type Permission struct {
	Base
	Modifier

	Title       string
	Key         PermissionKeyType
	Group       string
	Description string
}
