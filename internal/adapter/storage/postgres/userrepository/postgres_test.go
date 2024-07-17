package userrepository_test

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"path"
	"runtime"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db *sql.DB
	m  *migrate.Migrate
}

func (r *UserRepositoryTestSuite) SetupSuite() {
	var err error

	envPath, err := helper.GetEnvFilePath(".env.test")
	require.NoError(r.T(), err)

	configProvider := &config.Config{}
	conf := configProvider.GetConfig(envPath)

	require.NoError(r.T(), err)

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

	_, filename, _, _ := runtime.Caller(0)
	migrationDir := "file://" + path.Join(path.Dir(filename), "..", "migrations")

	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
	require.NoError(r.T(), err)

	r.m, err = migrate.NewWithDatabaseInstance(migrationDir, conf.DB.Name, driver)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), r.m)
	require.NoError(r.T(), r.m.Up())
}

func (r *UserRepositoryTestSuite) TearDownTest() {
	rows, err := r.db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
	require.NoError(r.T(), err)

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		require.NoError(r.T(), err)

		if tableName == "schema_migrations" {
			continue
		}

		queryTruncateTable := fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", tableName)
		_, err = r.db.Exec(queryTruncateTable)
		require.NoError(r.T(), err)
	}

	err = rows.Close()
	require.NoError(r.T(), err)
}

func (r *UserRepositoryTestSuite) TearDownSuite() {
	require.NoError(r.T(), r.m.Down())
}

func TestUserRepositoryTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping postgresql test in short mode")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}
