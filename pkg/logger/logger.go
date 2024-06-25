package logger

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

var once sync.Once
var zapSinLogger *zap.SugaredLogger

type Logger interface {
	Init(appName string)
	Debug(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	DebugF(template string, args ...interface{})
	Info(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	InfoF(template string, args ...interface{})
	Warn(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	WarnF(template string, args ...interface{})
	Error(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	ErrorF(template string, args ...interface{})
	Fatal(category Category, subCategory SubCategory, message string, extra map[ExtraKey]interface{})
	FatalF(template string, args ...interface{})
}

type zapLogger struct {
	config config.Log
	logger *zap.SugaredLogger
}

func (r *zapLogger) Init(appName string) {
	once.Do(func() {
		fileName := fmt.Sprintf(
			"%s%s.%s",
			r.config.FilePath,
			time.Now().Format("2006-01-02"),
			"log",
		)
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    r.config.MaxSize,
			MaxAge:     r.config.MaxAge,
			LocalTime:  true,
			MaxBackups: r.config.MaxBackups,
			Compress:   true,
		})

		stdoutWriter := zapcore.Lock(os.Stdout)

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			r.getLogLevel(),
		)

		stdoutCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			stdoutWriter,
			r.getLogLevel(),
		)

		teeCore := zapcore.NewTee(fileCore, stdoutCore)

		logger := zap.New(
			teeCore,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		).Sugar()
		zapSinLogger = logger.With("AppName", appName).With("LoggerName", "ZapLog")
	})
	r.logger = zapSinLogger
}

func NewLogger(appName string, config config.Log) Logger {
	logger := &zapLogger{config: config}
	logger.Init(appName)
	return logger
}

var zapLogLevelMapping = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

func (r *zapLogger) getLogLevel() zapcore.Level {
	level, exists := zapLogLevelMapping[r.config.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func prepareLogKeys(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["category"] = cat
	extra["subCategory"] = sub
	return MapToZapParams(extra)
}

func (r *zapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	r.logger.Debugw(msg, params...)
}

func (r *zapLogger) DebugF(template string, args ...interface{}) {
	r.logger.Debugf(template, args...)
}

func (r *zapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	r.logger.Infow(msg, params...)
}

func (r *zapLogger) InfoF(template string, args ...interface{}) {
	r.logger.Infof(template, args...)
}

func (r *zapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	r.logger.Warnw(msg, params...)
}

func (r *zapLogger) WarnF(template string, args ...interface{}) {
	r.logger.Warnf(template, args...)
}

func (r *zapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	r.logger.Errorw(msg, params...)
}

func (r *zapLogger) ErrorF(template string, args ...interface{}) {
	r.logger.Errorf(template, args...)
}

func (r *zapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogKeys(cat, sub, extra)
	r.logger.Fatalw(msg, params...)
}

func (r *zapLogger) FatalF(template string, args ...interface{}) {
	r.logger.Fatalf(template, args...)
}
