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

type PermissionRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewPermissionRepository(log logger.Logger, db *sql.DB) *PermissionRepository {
	return &PermissionRepository{
		log: log,
		db:  db,
	}
}

func (r PermissionRepository) GetUserPermissionKeys(ctx context.Context, userID uint64) ([]domain.PermissionKeyType, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT DISTINCT p.key
        FROM permissions p
        JOIN role_permissions rp ON p.id = rp.permission_id
        JOIN roles r ON rp.role_id = r.id AND r.deleted_at IS NULL
        JOIN access_controls ac on r.id = ac.role_id OR p.id = ac.permission_id AND ac.deleted_at IS NULL
        WHERE p.deleted_at IS NULL AND ac.user_id = $1`,
		userID,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}
	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}(rows)

	var permissionKeys []domain.PermissionKeyType
	var key domain.PermissionKeyType

	for rows.Next() {
		if err = rows.Scan(&key); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Success").Inc()

				r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
				return nil, serviceerror.New(serviceerror.RecordNotFound)
			}
			metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}
		permissionKeys = append(permissionKeys, key)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Success").Inc()

	return permissionKeys, nil
}

func (r PermissionRepository) List(ctx context.Context) ([]*domain.Permission, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT uuid, title, description, "group" FROM permissions WHERE deleted_at IS NULL;`)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}(rows)

	var permissions []*domain.Permission
	for rows.Next() {
		var permission domain.Permission
		if err = rows.Scan(&permission.UUID, &permission.Title, &permission.Description, &permission.Group); err != nil {
			metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}
		permissions = append(permissions, &permission)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "List", "Success").Inc()
	return permissions, nil
}
