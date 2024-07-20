package authrepository

import (
	"database/sql"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type ACLRepository struct {
	log logger.Logger
	tx  *sql.Tx
}

func NewACLRepository(log logger.Logger, tx *sql.Tx) *ACLRepository {
	return &ACLRepository{
		log: log,
		tx:  tx,
	}
}

func (r *ACLRepository) AssignRolesToUser(userID uint64, roleIDs []uint64) error {
	_, err := r.tx.Exec(`DELETE FROM access_controls WHERE user_id = $1`, userID)
	if err != nil {
		metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabasePrepare, err.Error(), nil)

		return serviceerror.NewServerError()
	}

	stmt, err := r.tx.Prepare(`INSERT INTO access_controls (user_id, role_id) VALUES ($1, $2)`)
	if err != nil {
		metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabasePrepare, err.Error(), nil)
		return serviceerror.NewServerError()
	}
	defer func(stmt *sql.Stmt) {
		if err = stmt.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), nil)
		}
	}(stmt)

	for _, roleID := range roleIDs {
		if _, err = stmt.Exec(userID, roleID); err != nil {
			metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
				"userID": userID,
				"roleID": roleID,
			})
			return serviceerror.NewServerError()
		}
	}

	metrics.DbCall.WithLabelValues("access_controls", "AssignRolesToUser", "Success").Inc()

	return nil
}
