package presenter_test

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPreparePermission(t *testing.T) {
	tests := []struct {
		name           string
		permission     *domain.Permission
		expectedResult *presenter.Permission
	}{
		{
			name:           "Nil Permission",
			permission:     nil,
			expectedResult: nil,
		},
		{
			name: "Valid Permission",
			permission: &domain.Permission{
				Base: domain.Base{
					UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
				},
				Title:       helper.StringPtr("Admin"),
				Group:       helper.StringPtr("user_role"),
				Description: helper.StringPtr("Admin description"),
			},
			expectedResult: &presenter.Permission{
				ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
				Title:       helper.StringPtr("Admin"),
				Group:       helper.StringPtr("user_role"),
				Description: helper.StringPtr("Admin description"),
			},
		},
		{
			name: "Invalid Permission with uuid equal nil",
			permission: &domain.Permission{
				Title:       helper.StringPtr("Admin"),
				Group:       helper.StringPtr("user_role"),
				Description: helper.StringPtr("Admin description"),
			},
			expectedResult: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.PreparePermission(test.permission)
			require.Equal(t, test.expectedResult, result)
		})
	}
}

func TestToPermissionCollection(t *testing.T) {
	tests := []struct {
		name           string
		permissions    []*domain.Permission
		expectedResult []presenter.Permission
	}{
		{
			name:           "Empty Permissions",
			permissions:    []*domain.Permission{},
			expectedResult: nil,
		},
		{
			name: "Valid Permissions",
			permissions: []*domain.Permission{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
					},
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
				{
					Base: domain.Base{
						UUID: uuid.MustParse("4f6d1172-19f6-4c42-8b2d-d09e6a8275e4"),
					},
					Title:       helper.StringPtr("User"),
					Group:       helper.StringPtr("user_management"),
					Description: helper.StringPtr("User management description"),
				},
			},
			expectedResult: []presenter.Permission{
				{
					ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
				{
					ID:          "4f6d1172-19f6-4c42-8b2d-d09e6a8275e4",
					Title:       helper.StringPtr("User"),
					Group:       helper.StringPtr("user_management"),
					Description: helper.StringPtr("User management description"),
				},
			},
		},
		{
			name: "Valid Permissions with One Empty Entry",
			permissions: []*domain.Permission{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
					},
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
				{}, // This is an empty entry
				{
					Base: domain.Base{
						UUID: uuid.MustParse("4f6d1172-19f6-4c42-8b2d-d09e6a8275e4"),
					},
					Title:       helper.StringPtr("User"),
					Group:       helper.StringPtr("user_management"),
					Description: helper.StringPtr("User management description"),
				},
			},
			expectedResult: []presenter.Permission{
				{
					ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
				{
					ID:          "4f6d1172-19f6-4c42-8b2d-d09e6a8275e4",
					Title:       helper.StringPtr("User"),
					Group:       helper.StringPtr("user_management"),
					Description: helper.StringPtr("User management description"),
				},
			},
		},
		{
			name: "Invalid Permissions with uuid equal nil",
			permissions: []*domain.Permission{
				{
					Base: domain.Base{
						UUID: uuid.MustParse("8f4a1582-6a67-4d85-950b-2d17049c7385"),
					},
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
				{
					Title:       helper.StringPtr("User"),
					Group:       helper.StringPtr("user_management"),
					Description: helper.StringPtr("User management description"),
				},
			},
			expectedResult: []presenter.Permission{
				{
					ID:          "8f4a1582-6a67-4d85-950b-2d17049c7385",
					Title:       helper.StringPtr("Admin"),
					Group:       helper.StringPtr("user_role"),
					Description: helper.StringPtr("Admin description"),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := presenter.ToPermissionCollection(test.permissions)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
