package presenter

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type Permission struct {
	ID          string  `json:"id,omitempty" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	Title       *string `json:"title,omitempty" example:"Admin"`
	Group       *string `json:"group,omitempty" example:"user_role"`
	Description *string `json:"description,omitempty" example:"Admin description"`
}

func PreparePermission(permission *domain.Permission) *Permission {
	if permission == nil || permission.UUID == uuid.Nil {
		return nil
	}

	return &Permission{
		ID:          permission.UUID.String(),
		Title:       permission.Title,
		Group:       permission.Group,
		Description: permission.Description,
	}
}

func ToPermissionCollection(permissions []*domain.Permission) []Permission {
	var response []Permission
	for _, permission := range permissions {
		result := PreparePermission(permission)
		if result != nil {
			response = append(response, *result)
		}
	}

	return response
}
