package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (r RoleRepository) Create(ctx context.Context, role domain.Role) error {
	res, err := r.db.ExecContext(
		ctx,
		`INSERT INTO roles (title, key, description) VALUES ($1, $2, $3)`,
		role.Title,
		role.Key,
		role.Description,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "Create", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
			logger.InsertDBArg: role,
		})
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("roles", "Create", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseInsert, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "Create", "Success").Inc()

	return nil
}

func (r RoleRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.db.QueryRowContext(ctx, "SELECT uuid, title, key, description FROM roles WHERE deleted_at IS NULL AND uuid = $1", uuid).
		Scan(&role.UUID, &role.Title, &role.Key, &role.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("roles", "GetByUUID", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.New(serviceerror.RecordNotFound)
		}
		metrics.DbCall.WithLabelValues("roles", "GetByUUID", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "GetByUUID", "Success").Inc()

	return &role, nil
}

func (r RoleRepository) List(ctx context.Context) ([]*domain.Role, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT uuid, title, key, description, is_default FROM roles WHERE deleted_at IS NULL")
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}(rows)

	var roles []*domain.Role

	for rows.Next() {

		var role domain.Role
		if err = rows.Scan(&role.UUID, &role.Title, &role.Key, &role.Description, &role.IsDefault); err != nil {
			metrics.DbCall.WithLabelValues("roles", "List", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, err
		}

		roles = append(roles, &role)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("roles", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, err
	}

	metrics.DbCall.WithLabelValues("roles", "List", "Success").Inc()
	return roles, nil
}

func (r RoleRepository) Update(ctx context.Context, role domain.Role, uuid uuid.UUID) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE roles SET title = $1, key = $2, description = $3 WHERE deleted_at IS NULL AND uuid = $4;",
		role.Title,
		role.Key,
		role.Description,
		uuid,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "Delete", "Success").Inc()

	return nil
}

func (r RoleRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE roles SET deleted_at = now(), key = $1 WHERE is_default = TRUE AND uuid = $2",
		uuid,
		uuid,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		if err != nil {
			metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
			return serviceerror.NewServerError()
		}

		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.New(serviceerror.IsNotDeletable)
	}

	metrics.DbCall.WithLabelValues("roles", "Delete", "Success").Inc()

	return nil
}

func (r RoleRepository) ExistKey(ctx context.Context, key string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM roles WHERE deleted_at IS NULL AND key = $1", key).
		Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("roles", "ExistKey", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return true, nil
		}
		metrics.DbCall.WithLabelValues("roles", "ExistKey", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return false, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "ExistKey", "Success").Inc()

	return count == 0, nil
}

func (r RoleRepository) GetRoleUser(ctx context.Context) (role domain.Role, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT id FROM roles WHERE key=$1`, domain.RoleKeyUser).
		Scan(&role.ID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("roles", "GetRoleUser", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return role, serviceerror.New(serviceerror.RecordNotFound)
		}
		metrics.DbCall.WithLabelValues("roles", "GetRoleUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return role, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "GetRoleUser", "Success").Inc()

	return role, nil
}
