package logger_test

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

const filePath = "./logs"

func teardown() {
	_ = os.RemoveAll(filePath)
}

func TestZapLoggerInitialization(t *testing.T) {
	defer teardown()

	cfg := config.Log{
		FilePath:   fmt.Sprintf("%s/logs", filePath),
		MaxSize:    1, // 1 MB
		MaxAge:     7, // 7 days
		MaxBackups: 3, // 3 backups
		Level:      "debug",
	}

	log := logger.NewLogger("TestApp", cfg)
	require.NotNil(t, log)
}

func TestZapLoggerMethods(t *testing.T) {
	defer teardown()

	cfg := config.Log{
		FilePath:   fmt.Sprintf("%s/logs", filePath),
		MaxSize:    1, // 1 MB
		MaxAge:     7, // 7 days
		MaxBackups: 3, // 3 backups
		Level:      "debug",
	}

	log := logger.NewLogger("TestApp", cfg)
	require.NotNil(t, log)

	extra := map[logger.ExtraKey]interface{}{
		"key1": "value1",
	}

	t.Run("Debug", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.Debug("category", "subCategory", "Debug message", extra)
		})
	})

	t.Run("Info", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.Info("category", "subCategory", "Info message", extra)
		})
	})

	t.Run("Warn", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.Warn("category", "subCategory", "Warn message", extra)
		})
	})

	t.Run("Error", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.Error("category", "subCategory", "Error message", extra)
		})
	})

	// Note: We are not actually calling Fatal here to avoid exiting the test process.
	// t.Run("Fatal", func(t *testing.T) {
	// 	require.NotPanics(t, func() {
	// 		defer func() {
	// 			if r := recover(); r != nil {
	// 				t.Log("Recovered in Fatal", r)
	// 			}
	// 		}()
	// 		log.Fatal("category", "subCategory", "Fatal message", extra)
	// 	})
	// })
}

func TestZapLoggerFormatMethods(t *testing.T) {
	defer teardown()

	cfg := config.Log{
		FilePath:   fmt.Sprintf("%s/logs", filePath),
		MaxSize:    1, // 1 MB
		MaxAge:     7, // 7 days
		MaxBackups: 3, // 3 backups
		Level:      "debug",
	}

	log := logger.NewLogger("TestApp", cfg)
	require.NotNil(t, log)

	t.Run("DebugF", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.DebugF("Debug %s", "message")
		})
	})

	t.Run("InfoF", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.InfoF("Info %s", "message")
		})
	})

	t.Run("WarnF", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.WarnF("Warn %s", "message")
		})
	})

	t.Run("ErrorF", func(t *testing.T) {
		require.NotPanics(t, func() {
			log.ErrorF("Error %s", "message")
		})
	})

	// Note: We are not actually calling FatalF here to avoid exiting the test process.
	// t.Run("FatalF", func(t *testing.T) {
	// 	require.NotPanics(t, func() {
	// 		defer func() {
	// 			if r := recover(); r != nil {
	// 				t.Log("Recovered in FatalF", r)
	// 			}
	// 		}()
	// 		log.FatalF("Fatal %s", "message")
	// 	})
	// })
}

func TestZapLoggerGetLogLevel(t *testing.T) {
	t.Run("Valid log level", func(t *testing.T) {
		cfg := config.Log{
			Level: "info",
		}
		log := logger.NewLogger("TestApp", cfg)
		zapLogger := log.(*logger.ZapLogger)
		require.Equal(t, zapcore.InfoLevel, zapLogger.GetLogLevel())
	})

	t.Run("Invalid log level defaults to debug", func(t *testing.T) {
		cfg := config.Log{
			Level: "invalid",
		}
		log := logger.NewLogger("TestApp", cfg)
		zapLogger := log.(*logger.ZapLogger)
		require.Equal(t, zapcore.DebugLevel, zapLogger.GetLogLevel())
	})
}

func TestPrepareLogKeys(t *testing.T) {
	t.Run("Non-nil extra map", func(t *testing.T) {
		extra := map[logger.ExtraKey]interface{}{
			"key1": "value1",
		}
		params := logger.PrepareLogKeys("category", "subCategory", extra)
		require.Contains(t, params, "category")
		require.Contains(t, params, "subCategory")
		require.Contains(t, params, "key1")
	})

	t.Run("Nil extra map", func(t *testing.T) {
		params := logger.PrepareLogKeys("category", "subCategory", nil)
		require.Contains(t, params, "category")
		require.Contains(t, params, "subCategory")
	})
}
