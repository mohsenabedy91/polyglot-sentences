package repository

import (
	"context"
	"database/sql"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type ACLRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewACLRepository(log logger.Logger, db *sql.DB) *ACLRepository {
	return &ACLRepository{
		log: log,
		db:  db,
	}
}

func (r ACLRepository) AssignRolesToUser(ctx context.Context, userID uint64, roleIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseBeginTransaction, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM access_controls WHERE user_id = $1`, userID)
	if err != nil {
		metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabasePrepare, err.Error(), nil)

		if err = tx.Rollback(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseRollback, err.Error(), nil)
		}

		return serviceerror.NewServerError()
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO access_controls (user_id, role_id) VALUES ($1, $2)`)
	if err != nil {
		metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabasePrepare, err.Error(), nil)

		if err = tx.Rollback(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseRollback, err.Error(), nil)
		}

		return serviceerror.NewServerError()
	}
	defer func(stmt *sql.Stmt) {
		if err = stmt.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), nil)
		}
	}(stmt)

	for _, roleID := range roleIDs {
		if _, err = stmt.ExecContext(ctx, userID, roleID); err != nil {
			metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
				"userID": userID,
				"roleID": roleID,
			})

			if err = tx.Rollback(); err != nil {
				r.log.Error(logger.Database, logger.DatabaseRollback, err.Error(), nil)
			}

			return serviceerror.NewServerError()
		}
	}

	if err = tx.Commit(); err != nil {
		metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseCommit, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Success").Inc()

	return nil
}
