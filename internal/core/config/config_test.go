package config_test

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func setup(t *testing.T) {
	err := os.MkdirAll("testdata/config", os.ModePerm)
	require.NoError(t, err)

	file, err := os.Create("testdata/config/.env")
	require.NoError(t, err)

	write, err := file.WriteString(`APP_NAME="MyApp" APP_ENV="development" APP_DEBUG="true" APP_GRACEFULLY_SHUTDOWN="10"`)
	require.NoError(t, err)
	require.Greater(t, write, 0)

	err = file.Close()
	require.NoError(t, err)

	config.ResetConfig()
}

func teardown(t *testing.T) {
	err := os.RemoveAll("testdata")
	require.NoError(t, err)
}

func TestLoadConfig(t *testing.T) {
	setup(t)
	defer teardown(t)

	conf := config.Config{}

	cfg, err := conf.LoadConfig("testdata/config/.env")
	require.NoError(t, err)

	require.Equal(t, "MyApp", cfg.App.Name)
	require.Equal(t, "development", cfg.App.Env)
	require.Equal(t, true, cfg.App.Debug)
	require.Equal(t, time.Duration(10), cfg.App.GracefullyShutdown)
}

func TestGetConfig(t *testing.T) {
	setup(t)
	defer teardown(t)

	conf := config.Config{}

	cfg := conf.GetConfig("testdata/config/.env")
	require.NotNil(t, cfg)
}

func TestGetConfigPanic(t *testing.T) {
	setup(t)
	defer teardown(t)

	conf := config.Config{}

	require.Panics(t, func() {
		_ = conf.GetConfig("testdata/config/.invalid")
	}, "Expected GetConfig to panic due to no such file or directory")
}
