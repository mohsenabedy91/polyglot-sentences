package port

import "context"

type UnitOfWork interface {
	BeginTx(ctx context.Context) error
	Commit() error
	Rollback() error
}

type AuthUnitOfWork interface {
	UnitOfWork

	RoleRepository() RoleRepository
	PermissionRepository() PermissionRepository
	ACLRepository() ACLRepository
	// Add other repositories as needed
}

type UserUnitOfWork interface {
	UnitOfWork

	UserRepository() UserRepository
	// Add other repositories as needed
}
