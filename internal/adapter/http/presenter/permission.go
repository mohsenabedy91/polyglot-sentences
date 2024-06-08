package presenter

import "github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"

type Permission struct {
	ID          string `json:"id" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	Title       string `json:"title" example:"Admin"`
	Group       string `json:"group" example:"user_role"`
	Description string `json:"description" example:"Admin description"`
}

func preparePermission(permission *domain.Permission) Permission {
	return Permission{
		ID:          permission.UUID.String(),
		Title:       permission.Title,
		Group:       permission.Group,
		Description: permission.Description,
	}
}

func ToPermissionCollection(permissions []*domain.Permission) []Permission {
	result := make([]Permission, len(permissions))
	for index, permission := range permissions {
		result[index] = preparePermission(permission)
	}

	return result
}
