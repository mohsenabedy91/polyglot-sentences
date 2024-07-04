package presenter_test

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrepareRole(t *testing.T) {
	tests := []struct {
		name           string
		role           *domain.Role
		expectedResult *presenter.Role
	}{
		{
			name:           "Nil Role",
			role:           nil,
			expectedResult: nil,
		},
		{
			name: "Valid Role",
			role: &domain.Role{
				Base: domain.Base{
					UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
				},
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []*domain.Permission{
					{
						Base: domain.Base{
							UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
						},
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
			expectedResult: &presenter.Role{
				ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []presenter.Permission{
					{
						ID:          "ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44",
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
		},
		{
			name: "Invalid Role with uuid equal nil",
			role: &domain.Role{
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []*domain.Permission{
					{
						Base: domain.Base{
							UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
						},
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
			expectedResult: nil,
		},
		{
			name: "Valid Role with Permission uuid equal nil",
			role: &domain.Role{
				Base: domain.Base{
					UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
				},
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []*domain.Permission{
					{
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
			expectedResult: &presenter.Role{
				ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.PrepareRole(test.role)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToRoleResource(t *testing.T) {
	tests := []struct {
		name           string
		role           *domain.Role
		expectedResult *presenter.Role
	}{
		{
			name:           "Nil Role",
			role:           nil,
			expectedResult: nil,
		},
		{
			name: "Valid Role",
			role: &domain.Role{
				Base: domain.Base{
					UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
				},
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []*domain.Permission{
					{
						Base: domain.Base{
							UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
						},
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
			expectedResult: &presenter.Role{
				ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []presenter.Permission{
					{
						ID:          "ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44",
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
		},
		{
			name: "Valid Role with One Of Permission Empty Entry",
			role: &domain.Role{
				Base: domain.Base{
					UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
				},
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []*domain.Permission{
					{}, // This is an empty entry
					{
						Base: domain.Base{
							UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
						},
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
			expectedResult: &presenter.Role{
				ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
				Title:       "Admin",
				Description: "Admin description",
				IsDefault:   true,
				Permissions: []presenter.Permission{
					{
						ID:          "ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44",
						Title:       helper.StringPtr("Read"),
						Group:       helper.StringPtr("user"),
						Description: helper.StringPtr("Read permission"),
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToRoleResource(test.role)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToRoleCollection(t *testing.T) {
	tests := []struct {
		name           string
		roles          []*domain.Role
		expectedResult []presenter.Role
	}{
		{
			name:           "Empty Roles",
			roles:          []*domain.Role{},
			expectedResult: nil,
		},
		{
			name: "Valid Roles",
			roles: []*domain.Role{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
					},
					Title:       "Admin",
					Description: "Admin description",
					IsDefault:   true,
					Permissions: []*domain.Permission{
						{
							Base: domain.Base{
								UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
							},
							Title:       helper.StringPtr("Read"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Read permission"),
						},
					},
				},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("4f6d1172-19f6-4c42-8b2d-d09e6a8275e4"),
					},
					Title:       "User",
					Description: "User description",
					IsDefault:   false,
					Permissions: []*domain.Permission{
						{
							Base: domain.Base{
								UUID: uuid.MustParse("e243264e-6e59-4363-8491-b9dedea6021c"),
							},
							Title:       helper.StringPtr("Write"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Write permission"),
						},
					},
				},
			},
			expectedResult: []presenter.Role{
				{
					ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
					Title:       "Admin",
					Description: "Admin description",
					IsDefault:   true,
					Permissions: []presenter.Permission{
						{
							ID:          "ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44",
							Title:       helper.StringPtr("Read"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Read permission"),
						},
					},
				},
				{
					ID:          "4f6d1172-19f6-4c42-8b2d-d09e6a8275e4",
					Title:       "User",
					Description: "User description",
					IsDefault:   false,
					Permissions: []presenter.Permission{
						{
							ID:          "e243264e-6e59-4363-8491-b9dedea6021c",
							Title:       helper.StringPtr("Write"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Write permission"),
						},
					},
				},
			},
		},
		{
			name: "Valid Roles with One Empty Entry",
			roles: []*domain.Role{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
					},
					Title:       "Admin",
					Description: "Admin description",
					IsDefault:   true,
					Permissions: []*domain.Permission{
						{}, // This is an empty entry
						{
							Base: domain.Base{
								UUID: uuid.MustParse("ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44"),
							},
							Title:       helper.StringPtr("Read"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Read permission"),
						},
					},
				},
				{}, // This is an empty entry
				{
					Base: domain.Base{
						UUID: uuid.MustParse("4f6d1172-19f6-4c42-8b2d-d09e6a8275e4"),
					},
					Title:       "User",
					Description: "User description",
					IsDefault:   false,
					Permissions: []*domain.Permission{
						{
							Base: domain.Base{
								UUID: uuid.MustParse("e243264e-6e59-4363-8491-b9dedea6021c"),
							},
							Title:       helper.StringPtr("Write"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Write permission"),
						},
					},
				},
			},
			expectedResult: []presenter.Role{
				{
					ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
					Title:       "Admin",
					Description: "Admin description",
					IsDefault:   true,
					Permissions: []presenter.Permission{
						{
							ID:          "ef9b70b8-aca5-4a78-bdd7-ee41f7aabc44",
							Title:       helper.StringPtr("Read"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Read permission"),
						},
					},
				},
				{
					ID:          "4f6d1172-19f6-4c42-8b2d-d09e6a8275e4",
					Title:       "User",
					Description: "User description",
					IsDefault:   false,
					Permissions: []presenter.Permission{
						{
							ID:          "e243264e-6e59-4363-8491-b9dedea6021c",
							Title:       helper.StringPtr("Write"),
							Group:       helper.StringPtr("user"),
							Description: helper.StringPtr("Write permission"),
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToRoleCollection(test.roles)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
