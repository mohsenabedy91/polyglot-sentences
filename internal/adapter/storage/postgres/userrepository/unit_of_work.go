package userrepository

import (
	"context"
	"database/sql"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type UnitOfWork interface {
	BeginTx(ctx context.Context) error
	Commit() error
	Rollback() error

	UserRepository() *UserRepository
	// Add other repositories as needed
}

type unitOfWork struct {
	log logger.Logger
	db  *sql.DB
	tx  *sql.Tx

	userRepository *UserRepository
	// Add other repositories as needed
}

func NewUnitOfWork(log logger.Logger, db *sql.DB) UnitOfWork {
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
	r.userRepository = NewUserRepository(r.log, tx)
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

func (r *unitOfWork) UserRepository() *UserRepository {
	return r.userRepository
}
