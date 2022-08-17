package logger

import (
	"fmt"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
)

type Logger struct {
	s *zap.SugaredLogger
}

func New(level config.LogLevel) (*Logger, error) {
	l := zap.NewAtomicLevel()
	switch level {
	case config.LogLevelDebug:
		l.SetLevel(zap.DebugLevel)
	case config.LogLevelError:
		l.SetLevel(zap.ErrorLevel)
	case config.LogLevelWarn:
		l.SetLevel(zap.WarnLevel)
	case config.LogLevelInfo:
		l.SetLevel(zap.InfoLevel)
	default:
		return nil, fmt.Errorf("invalid log level: %s", level)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = l
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	sugar := logger.WithOptions(zap.AddCallerSkip(2)).Sugar()
	return &Logger{
		s: sugar,
	}, nil
}

func (l Logger) Info(data interface{}) {
	l.s.Info(data)
}

func (l Logger) Error(data interface{}) {
	l.s.Error(data)
}

func (l Logger) Debug(data interface{}) {
	l.s.Debug(data)
}

func (l Logger) Warn(data interface{}) {
	l.s.Warn(data)
}
