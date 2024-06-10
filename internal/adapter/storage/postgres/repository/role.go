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
	err := r.db.QueryRowContext(ctx, "SELECT uuid, title, key, description, is_default FROM roles WHERE deleted_at IS NULL AND uuid = $1", uuid).
		Scan(&role.UUID, &role.Title, &role.Key, &role.Description, &role.IsDefault)
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

		r.log.Error(logger.Database, logger.DatabaseDelete, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		if err != nil {
			metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseDelete, err.Error(), nil)
			return serviceerror.NewServerError()
		}

		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseDelete, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
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

func (r RoleRepository) GetPermissions(ctx context.Context, roleUUID uuid.UUID) (*domain.Role, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT r.uuid, r.title, r.description, p.uuid, p.title, p.group, p.description
				FROM roles AS r
				LEFT JOIN role_permissions AS rp ON rp.role_id = r.id
				LEFT JOIN permissions AS p ON rp.permission_id = p.id
				WHERE r.uuid = $1`,
		roleUUID,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "GetPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}(rows)

	var role domain.Role
	var permissions []*domain.Permission

	for rows.Next() {
		var permission domain.Permission
		err = rows.Scan(
			&role.UUID,
			&role.Title,
			&role.Description,
			&permission.UUID,
			&permission.Title,
			&permission.Group,
			&permission.Description,
		)
		if err != nil {
			metrics.DbCall.WithLabelValues("roles", "GetPermissions", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}

		if permission.UUID == uuid.Nil {
			continue
		}

		permissions = append(permissions, &permission)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("roles", "GetPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	if role.UUID == uuid.Nil {
		metrics.DbCall.WithLabelValues("roles", "GetPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, fmt.Sprintf("There is any role for %s", roleUUID.String()), nil)
		return nil, serviceerror.New(serviceerror.RecordNotFound)
	}

	metrics.DbCall.WithLabelValues("roles", "GetPermissions", "Success").Inc()

	role.Permissions = permissions
	return &role, nil
}

func (r RoleRepository) GetRoleKeys(ctx context.Context, userID uint64) ([]domain.RoleKeyType, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT DISTINCT r.key FROM access_controls AS ac
         		INNER JOIN roles r on r.id = ac.role_id AND r.deleted_at IS NULL
				WHERE ac.deleted_at IS NULL AND ac.user_id = $1`,
		userID,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "GetRoleKeys", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func() {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}()

	var keys []domain.RoleKeyType
	for rows.Next() {
		var key domain.RoleKeyType
		if err = rows.Scan(&key); err != nil {
			metrics.DbCall.WithLabelValues("roles", "GetRoleKeys", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}
		keys = append(keys, key)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("roles", "GetRoleKeys", "Failed").Inc()
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)

		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("roles", "GetRoleKeys", "Success").Inc()

	return keys, nil
}
