package logger

import "go.uber.org/zap"

type ZapLogger struct {
	zapSugarLogger *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &ZapLogger{zapSugarLogger: logger.Sugar()}
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
