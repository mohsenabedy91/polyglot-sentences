package authrepository

import (
	"context"
	"database/sql"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type unitOfWork struct {
	log logger.Logger
	db  *sql.DB
	tx  *sql.Tx

	roleRepository       port.RoleRepository
	permissionRepository port.PermissionRepository
	aclRepository        port.ACLRepository
	// Add other repositories as needed
}

func NewUnitOfWork(log logger.Logger, db *sql.DB) port.AuthUnitOfWork {
	return &unitOfWork{
		log: log,
		db:  db,
	}
}

func (r *unitOfWork) BeginTx(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseBeginTransaction, err.Error(), nil)

		return serviceerror.NewServerError()
	}

	r.tx = tx

	r.roleRepository = NewRoleRepository(r.log, tx)
	r.permissionRepository = NewPermissionRepository(r.log, tx)
	r.aclRepository = NewACLRepository(r.log, tx)
	// Initialize other repositories as needed

	return nil
}

func (r *unitOfWork) Commit() error {

	if err := r.tx.Commit(); err != nil {
		r.log.Error(logger.Database, logger.DatabaseCommit, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	return nil
}

func (r *unitOfWork) Rollback() error {

	if err := r.tx.Rollback(); err != nil {
		r.log.Error(logger.Database, logger.DatabaseRollback, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	return nil
}

func (r *unitOfWork) RoleRepository() port.RoleRepository {
	return r.roleRepository
}

func (r *unitOfWork) PermissionRepository() port.PermissionRepository {
	return r.permissionRepository
}

func (r *unitOfWork) ACLRepository() port.ACLRepository {
	return r.aclRepository
}
