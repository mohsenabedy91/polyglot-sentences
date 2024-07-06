package userrepository

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

	userRepository port.UserRepository
	// Add other repositories as needed
}

func NewUnitOfWork(log logger.Logger, db *sql.DB) port.UserUnitOfWork {
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

func (r *unitOfWork) UserRepository() port.UserRepository {
	return r.userRepository
}
