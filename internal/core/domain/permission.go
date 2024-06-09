package domain

type PermissionKeyType string

const (
	PermissionKeyCreateUser              PermissionKeyType = "CREATE_USER"
	PermissionKeyReadUser                PermissionKeyType = "READ_USER"
	PermissionKeyUpdateUser              PermissionKeyType = "UPDATE_USER"
	PermissionKeyDeleteUser              PermissionKeyType = "DELETE_USER"
	PermissionKeyCreateRole              PermissionKeyType = "CREATE_ROLE"
	PermissionKeyReadRole                PermissionKeyType = "READ_ROLE"
	PermissionKeyUpdateRole              PermissionKeyType = "UPDATE_ROLE"
	PermissionKeyDeleteRole              PermissionKeyType = "DELETE_ROLE"
	PermissionKeyReadPermission          PermissionKeyType = "READ_PERMISSION"
	PermissionKeyAssignRolesToUser       PermissionKeyType = "ASSIGN_ROLES_TO_USER"
	PermissionKeyReadUserRoles           PermissionKeyType = "READ_USER_ROLES"
	PermissionKeyAssignPermissionsToRole PermissionKeyType = "ASSIGN_PERMISSIONS_TO_ROLE"
	PermissionKeyReadRolePermissions     PermissionKeyType = "READ_ROLE_PERMISSIONS"
)

type Permission struct {
	Base
	Modifier

	Title       *string
	Key         *PermissionKeyType
	Group       *string
	Description *string
}
