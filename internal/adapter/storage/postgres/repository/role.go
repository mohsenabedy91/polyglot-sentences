package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type RoleRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewRoleRepository(log logger.Logger, db *sql.DB) *RoleRepository {
	return &RoleRepository{
		log: log,
		db:  db,
	}
}

func (r RoleRepository) GetUserRole(ctx context.Context) (role domain.Role, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT id FROM roles WHERE key=$1`, domain.RoleKeyUser).
		Scan(&role.ID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("users", "GetUserRole", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return role, serviceerror.New(serviceerror.RecordNotFound)
		}
		metrics.DbCall.WithLabelValues("users", "GetUserRole", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return role, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "GetUserRole", "Success").Inc()

	return role, nil
}
