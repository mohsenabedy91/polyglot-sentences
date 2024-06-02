package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (r ACLRepository) AddUserRole(ctx context.Context, userID, roleID uint64) error {
	res, err := r.db.ExecContext(ctx, `INSERT INTO access_controls (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "AddUserRole", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
			logger.InsertDBArg: userID,
		})
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("users", "AddUserRole", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseInsert, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "AddUserRole", "Success").Inc()

	return nil
}
