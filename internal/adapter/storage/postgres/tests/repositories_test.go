package tests

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"path"
	"runtime"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestSuite struct {
	suite.Suite
	db *sql.DB
	tx *sql.Tx
	m  *migrate.Migrate
}

func (r *TestSuite) SetupSuite() {
	var err error

	envPath, err := helper.GetEnvFilePath(".env.test")
	require.NoError(r.T(), err)

	configProvider := &config.Config{}
	conf := configProvider.GetConfig(envPath)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Name,
		conf.DB.Postgres.SSLMode,
		conf.DB.Postgres.Timezone,
	)

	r.db, err = sql.Open("postgres", dsn)
	require.NoError(r.T(), err)

	_, filename, _, ok := runtime.Caller(0)
	require.True(r.T(), ok)

	migrationDir := "file://" + path.Join(path.Dir(filename), "..", "migrations")

	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
	require.NoError(r.T(), err)

	r.m, err = migrate.NewWithDatabaseInstance(migrationDir, conf.DB.Name, driver)
	require.NoError(r.T(), err)
	require.NotNil(r.T(), r.m)

	_ = r.m.Down()
	require.NoError(r.T(), r.m.Up())
}

func (r *TestSuite) SetupTest() {
	var err error
	r.tx, err = r.db.Begin()
	require.NoError(r.T(), err)
}

func (r *TestSuite) TearDownTest() {
	require.NoError(r.T(), r.tx.Rollback())
}

func (r *TestSuite) TearDownSuite() {
	require.NoError(r.T(), r.m.Down())
	require.NoError(r.T(), r.db.Close())
}

func (r *TestSuite) GetTx() *sql.Tx {
	return r.tx
}

func TestRepositoriesTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping postgresql test in short mode")
	}

	t.Parallel()
	suite.Run(t, new(RoleRepositoryTestSuite))
	suite.Run(t, new(UserRepositoryTestSuite))
	suite.Run(t, new(PermissionRepositoryTestSuite))
}

func insertUser(t *testing.T, tx *sql.Tx, user *domain.User) *domain.User {
	require.NoError(t, tx.QueryRow(
		`INSERT INTO users (first_name, last_name, email, password, status, google_id, avatar, created_by) 
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
							RETURNING id, uuid`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Status,
		user.GoogleID,
		user.Avatar,
		user.Modifier.CreatedBy,
	).Scan(&user.Base.ID, &user.Base.UUID))

	return user
}

func insertRole(t *testing.T, tx *sql.Tx, role *domain.Role) *domain.Role {
	require.NoError(t, tx.QueryRow(
		"INSERT INTO roles (title, key, description, created_by) VALUES ($1, $2, $3, $4) RETURNING id, uuid",
		role.Title,
		role.Key,
		role.Description,
		role.Modifier.CreatedBy,
	).Scan(&role.Base.ID, &role.Base.UUID))

	return role
}

func insertPermission(t *testing.T, tx *sql.Tx, permission *domain.Permission) *domain.Permission {
	require.NoError(t, tx.QueryRow(
		"INSERT INTO permissions (title, key, description, created_by) VALUES ($1, $2, $3, $4) RETURNING id, uuid",
		permission.Title,
		permission.Key,
		permission.Description,
		permission.Modifier.CreatedBy,
	).Scan(&permission.Base.ID, &permission.Base.UUID))

	return permission
}

func insertRolePermission(t *testing.T, tx *sql.Tx, roleID uint64, permissionID uint64) {
	_, err := tx.Exec(
		"INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2);",
		roleID,
		permissionID,
	)
	require.NoError(t, err)
}

func addRoleToUser(t *testing.T, tx *sql.Tx, userID uint64, roleID uint64) {
	_, err := tx.Exec(
		"INSERT INTO access_controls (user_id, role_id) VALUES ($1, $2);",
		userID,
		roleID,
	)
	require.NoError(t, err)
}
