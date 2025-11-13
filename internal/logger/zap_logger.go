package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type ZapLogger struct {
	zapSugarLogger *zap.SugaredLogger
}

func NewZapLogger(isProduction bool, isStackTrace bool) *ZapLogger {
	var zapLogger *zap.Logger
	var err error

	var cfg zap.Config
	if isProduction {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.DisableStacktrace = !isStackTrace
	zapLogger, err = cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to create zap logger: %w", err))
	}
	return &ZapLogger{zapSugarLogger: zapLogger.Sugar()}
}

func (l *ZapLogger) Debug(msg string, args ...any) {
	l.zapSugarLogger.Debugw(msg, args...)
}

func (l *ZapLogger) Info(msg string, args ...any) {
	l.zapSugarLogger.Infow(msg, args...)
}

func (l *ZapLogger) Warn(msg string, args ...any) {
	l.zapSugarLogger.Warnw(msg, args...)
}

func (l *ZapLogger) Error(msg string, args ...any) {
	l.zapSugarLogger.Errorw(msg, args...)
}

func (l *ZapLogger) Fatal(msg string, args ...any) {
	l.zapSugarLogger.Fatalw(msg, args...)
}
