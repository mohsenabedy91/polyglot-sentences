package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
)

// UserRepository implements port.UserRepository interface and provides access to the postgres database
type UserRepository struct {
	log logger.Logger
	db  *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(log logger.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

func (r *UserRepository) IsEmailUnique(ctx context.Context, email string) (bool, serviceerror.Error) {
	email = strings.ToLower(email)
	var count int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND LOWER(email) = $1`,
		email,
	).Scan(&count)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}
		return false, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	return count == 0, nil
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) serviceerror.Error {
	res, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (first_name, last_name, email, password, status) VALUES ($1, $2, $3, $4, $5)",
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		domain.UserStatusUnVerified,
	)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
			logger.InsertDBArg: user,
		})
		return serviceerror.NewServiceError(serviceerror.ServerError)
	}

	affected, err := res.RowsAffected()
	if err != nil || affected <= 0 {
		if err != nil {
			r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), nil)
			return serviceerror.NewServiceError(serviceerror.ServerError)
		}
		r.log.Error(logger.Database, logger.DatabaseInsert, "There is any effected row in DB.", nil)
		return serviceerror.NewServiceError(serviceerror.ServerError)
	}

	return nil
}

func (r *UserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, serviceerror.Error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, uuid, first_name, last_name, email, status FROM users WHERE deleted_at IS NULL AND uuid = $1",
		uuid,
	)
	user, err := scanUser(row)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceerror.NewServiceError(serviceerror.RecordNotFound)
		}
		return nil, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	if !user.IsActive() {
		r.log.Warn(logger.Database, logger.DatabaseSelect, "The User is inactive", map[logger.ExtraKey]interface{}{
			logger.SelectDBArg: user,
		})
		return nil, serviceerror.NewServiceError(serviceerror.UserInActive)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, serviceerror.Error) {
	user := &domain.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT uuid, password FROM users WHERE deleted_at IS NULL AND status = $1 AND LOWER(email) = $2",
		domain.UserStatusActive,
		strings.ToLower(email),
	).Scan(&user.UUID, &user.Password)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceerror.NewServiceError(serviceerror.RecordNotFound)
		}
		return nil, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, serviceerror.Error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, uuid, first_name, last_name, email, status FROM users WHERE deleted_at IS NULL",
	)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return
		}
	}(rows)

	var users []domain.User

	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, serviceerror.NewServiceError(serviceerror.RecordNotFound)
			}
			return nil, serviceerror.NewServiceError(serviceerror.ServerError)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServiceError(serviceerror.ServerError)
	}

	return users, nil
}

func scanUser(scanner postgres.Scanner) (domain.User, error) {
	var user domain.User

	err := scanner.Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Status)

	return user, err
}
