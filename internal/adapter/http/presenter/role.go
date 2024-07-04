package presenter

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type Role struct {
	ID          string       `json:"id" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	Title       string       `json:"title" example:"Admin"`
	Description string       `json:"description,omitempty" example:"Admin description"`
	IsDefault   bool         `json:"isDefault,omitempty" example:"true"`
	Permissions []Permission `json:"permissions,omitempty"`
}

func PrepareRole(role *domain.Role) *Role {
	if role == nil || role.UUID == uuid.Nil {
		return nil
	}

	return &Role{
		ID:          role.UUID.String(),
		Title:       role.Title,
		Description: role.Description,
		IsDefault:   role.IsDefault,
		Permissions: ToPermissionCollection(role.Permissions),
	}
}

func ToRoleResource(role *domain.Role) *Role {
	return PrepareRole(role)
}

func ToRoleCollection(roles []*domain.Role) []Role {
	var response []Role
	for _, role := range roles {
		result := PrepareRole(role)
		if result != nil {
			response = append(response, *result)
		}
	}

	return response
}
