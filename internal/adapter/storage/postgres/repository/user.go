package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/metrics"
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

func (r *UserRepository) IsEmailUnique(ctx context.Context, email string) (bool, error) {
	email = strings.ToLower(email)
	var count int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND LOWER(email) = $1`,
		email,
	).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("users", "IsEmailUnique", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return true, nil
		}
		metrics.DbCall.WithLabelValues("users", "IsEmailUnique", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return false, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "IsEmailUnique", "Success").Inc()

	return count == 0, nil
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (first_name, last_name, email, password, status, google_id) 
							VALUES ($1, $2, $3, $4, $5, $6) 
							RETURNING id, uuid`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Status,
		user.GoogleID,
	).Scan(&user.ID, &user.UUID)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "Save", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseInsert, err.Error(), map[logger.ExtraKey]interface{}{
			logger.InsertDBArg: user,
		})
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "Save", "Success").Inc()
	return user, nil
}

func (r *UserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		"SELECT id, uuid, first_name, last_name, email, status FROM users WHERE deleted_at IS NULL AND uuid = $1",
		uuid,
	)
	user, err := scanUser(row)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "GetByUUID", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, serviceerror.New(serviceerror.RecordNotFound)
		}
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "GetByUUID", "Success").Inc()

	if !user.IsActive() {
		r.log.Warn(logger.Database, logger.DatabaseSelect, "The User is inactive", map[logger.ExtraKey]interface{}{
			logger.SelectDBArg: user,
		})
		return nil, serviceerror.New(serviceerror.UserInActive)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	var googleID sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, uuid, first_name, last_name, email, password, welcome_message_sent, google_id FROM users 
					WHERE deleted_at IS NULL AND status IN ($1, $2) AND LOWER(email) = $3`,
		domain.UserStatusUnverifiedStr,
		domain.UserStatusActive,
		strings.ToLower(email),
	).Scan(&user.ID, &user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.WelcomeMessageSent, &googleID)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			metrics.DbCall.WithLabelValues("users", "GetByEmail", "Success").Inc()

			r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, nil
		}
		metrics.DbCall.WithLabelValues("users", "GetByEmail", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	user.SetGoogleID(googleID)

	metrics.DbCall.WithLabelValues("users", "GetByEmail", "Success").Inc()

	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, uuid, first_name, last_name, email, status FROM users WHERE deleted_at IS NULL",
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return
		}
	}(rows)

	var users []domain.User

	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				metrics.DbCall.WithLabelValues("users", "List", "Success").Inc()

				r.log.Warn(logger.Database, logger.DatabaseSelect, err.Error(), nil)
				return nil, serviceerror.New(serviceerror.RecordNotFound)
			}
			metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

			r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
			return nil, serviceerror.NewServerError()
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		metrics.DbCall.WithLabelValues("users", "List", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	metrics.DbCall.WithLabelValues("users", "List", "Success").Inc()

	return users, nil
}

func (r *UserRepository) VerifiedEmail(ctx context.Context, email string) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE users SET email_verified_at = now(), status = $1, updated_at = NOW() WHERE deleted_at IS NULL AND email = $2",
		domain.UserStatusActive,
		email,
	)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "VerifiedEmail", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := res.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("users", "VerifiedEmail", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}
	metrics.DbCall.WithLabelValues("users", "VerifiedEmail", "Success").Inc()

	return nil
}

func (r *UserRepository) MarkWelcomeMessageSent(ctx context.Context, id uint64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET welcome_message_sent = true, updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "MarkWelcomeMessageSent", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := result.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("users", "MarkWelcomeMessageSent", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}
	metrics.DbCall.WithLabelValues("users", "MarkWelcomeMessageSent", "Success").Inc()

	return nil
}

func (r *UserRepository) UpdateGoogleID(ctx context.Context, ID uint64, googleID string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET google_id = $1, updated_at = NOW() WHERE id = $2;", googleID, ID)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "UpdateGoogleID", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := result.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("users", "UpdateGoogleID", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}
	metrics.DbCall.WithLabelValues("users", "UpdateGoogleID", "Success").Inc()

	return nil
}

func (r *UserRepository) UpdateLastLoginTime(ctx context.Context, ID uint64) error {
	result, err := r.db.ExecContext(ctx, "UPDATE users SET last_login = now(), updated_at = NOW() WHERE id = $1;", ID)
	if err != nil {
		metrics.DbCall.WithLabelValues("users", "UpdateLastLoginTime", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, err.Error(), nil)
		return serviceerror.NewServerError()
	}

	if affected, err := result.RowsAffected(); err != nil || affected <= 0 {
		metrics.DbCall.WithLabelValues("users", "UpdateLastLoginTime", "Failed").Inc()

		r.log.Error(logger.Database, logger.DatabaseUpdate, fmt.Sprintf("There is any effected row in DB: %v", err), nil)
		return serviceerror.NewServerError()
	}
	metrics.DbCall.WithLabelValues("users", "UpdateLastLoginTime", "Success").Inc()

	return nil
}

func scanUser(scanner postgres.Scanner) (domain.User, error) {
	var user domain.User
	var firstName sql.NullString
	var lastName sql.NullString

	err := scanner.Scan(&user.ID, &user.UUID, &firstName, &lastName, &user.Email, &user.Status)

	user.SetFirstName(firstName).SetLastName(lastName)

	return user, err
}
