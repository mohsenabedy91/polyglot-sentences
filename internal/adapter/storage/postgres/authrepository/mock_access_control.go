package authrepository

import "github.com/stretchr/testify/mock"

type MockACLRepository struct {
	mock.Mock
}

func (r *MockACLRepository) AssignRolesToUser(userID uint64, roleIDs []uint64) error {
	args := r.Called(userID, roleIDs)
	return args.Error(0)
}
