package presenter

import "github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"

type Role struct {
	ID          string `json:"id" example:"8f4a1582-6a67-4d85-950b-2d17049c7385"`
	Title       string `json:"title" example:"Admin"`
	Description string `json:"description" example:"Admin description"`
	IsDefault   bool   `json:"isDefault" example:"true"`
}

func prepareRole(role *domain.Role) Role {
	if role == nil {
		return Role{}
	}

	return Role{
		ID:          role.UUID.String(),
		Title:       role.Title,
		Description: role.Description,
		IsDefault:   role.IsDefault,
	}
}

func ToRoleResource(role *domain.Role) Role {
	return prepareRole(role)
}

func ToRoleCollection(roles []*domain.Role) []Role {
	result := make([]Role, len(roles))
	for index, role := range roles {
		result[index] = prepareRole(role)
	}

	return result
}
