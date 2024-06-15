package requests

type RoleUUIDUri struct {
	UUIDStr string `uri:"roleID" binding:"required,uuid" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
}

type RoleCreate struct {
	Title       string `json:"title" binding:"required,min=3,max=64,role_title" example:"admin"`
	Description string `json:"description" binding:"required,min=5" example:"admin access to user management, permission management, etc"`
}

type RoleUpdate struct {
	Title       string `json:"title" binding:"required,min=3,max=64,role_title" example:"admin"`
	Description string `json:"description" binding:"required,min=5" example:"admin access to user management, permission management, etc"`
}

type SyncPermissions struct {
	Permissions []string `json:"permissions" binding:"required,dive,uuid" example:"550e8400-e29b-41d4-a716-446655440000,550e8400-e29b-41d4-a716-446655440001"`
}
