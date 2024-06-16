package repository

import (
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
	tx  *sql.Tx
}

func NewRoleRepository(log logger.Logger, tx *sql.Tx) *RoleRepository {
	return &RoleRepository{
		log: log,
		tx:  tx,
	}
}

func (r *RoleRepository) Create(role domain.Role) error {
	res, err := r.tx.Exec(
		`INSERT INTO roles (title, key, description, created_by) VALUES ($1, $2, $3, $4)`,
		role.Title,
		role.Key,
		role.Description,
		role.CreatedBy,
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

func (r *RoleRepository) GetByUUID(uuid uuid.UUID) (*domain.Role, error) {
	var role domain.Role
	err := r.tx.QueryRow("SELECT id, uuid, title, key, description, is_default FROM roles WHERE deleted_at IS NULL AND uuid = $1", uuid).
		Scan(&role.ID, &role.UUID, &role.Title, &role.Key, &role.Description, &role.IsDefault)
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

func (r *RoleRepository) List() ([]*domain.Role, error) {
	rows, err := r.tx.Query("SELECT uuid, title, key, description, is_default FROM roles WHERE deleted_at IS NULL")
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

func (r *RoleRepository) Update(role domain.Role, uuid uuid.UUID) error {
	res, err := r.tx.Exec(
		"UPDATE roles SET title = $1, key = $2, description = $3, updated_at = now(), updated_by = $4 WHERE deleted_at IS NULL AND uuid = $5;",
		role.Title,
		role.Key,
		role.Description,
		role.UpdatedBy,
		uuid,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected < 0 {
		metrics.DbCall.WithLabelValues("roles", "Delete", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("%v", err), nil)
		return serviceerror.NewServerError()
	} else if affected == 0 {

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.New(serviceerror.NoRowsEffected)
	}

	metrics.DbCall.WithLabelValues("roles", "Delete", "Success").Inc()

	return nil
}

func (r *RoleRepository) Delete(uuid uuid.UUID, deletedBy uint64) error {
	res, err := r.tx.Exec(
		"UPDATE roles SET deleted_at = now(), key = $1, deleted_by = $2 WHERE is_default = FALSE AND uuid = $3;",
		uuid,
		deletedBy,
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

func (r *RoleRepository) ExistKey(key domain.RoleKeyType) (bool, error) {
	var count int
	err := r.tx.QueryRow("SELECT count(*) FROM roles WHERE deleted_at IS NULL AND key = $1", key).
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

func (r *RoleRepository) GetRoleUser() (role domain.Role, err error) {
	err = r.tx.QueryRow(`SELECT id FROM roles WHERE key=$1`, domain.RoleKeyUser).
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

func (r *RoleRepository) GetPermissions(roleUUID uuid.UUID) (*domain.Role, error) {
	rows, err := r.tx.Query(
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

func (r *RoleRepository) GetRoleKeys(userID uint64) ([]domain.RoleKeyType, error) {
	rows, err := r.tx.Query(
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

func (r *RoleRepository) SyncPermissions(roleID uint64, permissionIDs []uint64) error {
	result, err := r.tx.Exec("DELETE FROM role_permissions WHERE role_id = $1", roleID)
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "SyncPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseDelete, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := result.RowsAffected(); err != nil || affected < 0 {
		metrics.DbCall.WithLabelValues("roles", "SyncPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("%v", err), nil)
		return serviceerror.NewServerError()
	}

	stmt, err := r.tx.Prepare("INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2);")
	if err != nil {
		metrics.DbCall.WithLabelValues("roles", "SyncPermissions", "Failed").Inc()
		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	defer func(stmt *sql.Stmt) {
		if err = stmt.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), nil)
		}
	}(stmt)

	for _, permissionID := range permissionIDs {
		if _, err = stmt.Exec(roleID, permissionID); err != nil {
			metrics.DbCall.WithLabelValues("roles", "SyncPermissions", "Failed").Inc()
			r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
			return serviceerror.NewServerError()
		}
	}

	metrics.DbCall.WithLabelValues("roles", "SyncPermissions", "Success").Inc()

	return nil
}
