package authrepository

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
)

type PermissionRepository struct {
	log logger.Logger
	tx  *sql.Tx
}

func NewPermissionRepository(log logger.Logger, tx *sql.Tx) *PermissionRepository {
	return &PermissionRepository{
		log: log,
		tx:  tx,
	}
}

func (r *PermissionRepository) GetUserPermissionKeys(userID uint64) ([]domain.PermissionKeyType, error) {
	rows, err := r.tx.Query(
		`SELECT DISTINCT p.key FROM access_controls AS ac
				LEFT JOIN roles r on r.id = ac.role_id AND r.deleted_at IS NULL
				LEFT JOIN role_permissions rp on rp.role_id = r.id AND r.deleted_at IS NULL
				LEFT JOIN permissions AS p on (p.id = rp.permission_id OR p.id = ac.permission_id) AND p.deleted_at IS NULL
				WHERE ac.deleted_at IS NULL AND ac.user_id = $1`,
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
	var key *domain.PermissionKeyType

	for rows.Next() {
		if err = rows.Scan(&key); err != nil {
			metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}
		if key == nil {
			continue
		}
		permissionKeys = append(permissionKeys, *key)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "GetUserPermissionKeys", "Success").Inc()

	return permissionKeys, nil
}

func (r *PermissionRepository) List() ([]*domain.Permission, error) {
	rows, err := r.tx.Query(`SELECT id, uuid, title, description, "group" FROM permissions WHERE deleted_at IS NULL;`)
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
		if err = rows.Scan(&permission.Base.ID, &permission.Base.UUID, &permission.Title, &permission.Description, &permission.Group); err != nil {
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

func (r *PermissionRepository) FilterValidPermissions(uuids []uuid.UUID) ([]uint64, error) {
	placeholders := strings.Join(helper.MakeSQLPlaceholders(uint(len(uuids))), ",")
	args := make([]interface{}, len(uuids))
	for i, u := range uuids {
		args[i] = u.String()
	}
	query := `SELECT id FROM permissions WHERE deleted_at IS NULL AND uuid IN (` + placeholders + `);`
	rows, err := r.tx.Query(query, args...)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "FilterValidPermissions", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func() {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		}
	}()

	var validPermissions []uint64
	for rows.Next() {
		var permission uint64
		if err = rows.Scan(&permission); err != nil {
			metrics.DbCall.WithLabelValues("users", "FilterValidPermissions", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}

		validPermissions = append(validPermissions, permission)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("users", "FilterValidPermissions", "Failed").Inc()
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)

		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "FilterValidPermissions", "Success").Inc()

	return validPermissions, nil
}
