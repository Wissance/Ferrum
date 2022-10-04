package logging

import (
	"Ferrum/config"
	"fmt"
	"github.com/mattn/go-colorable"
	"io"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ttys3/rotatefilehook"
)

const (
	timestampFormat = time.RFC822
	defaultLogLevel = log.InfoLevel
)

var logLevels = map[string]log.Level{
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"debug": log.DebugLevel,
	"trace": log.TraceLevel,
}

type AppLogger struct {
	logger    *log.Logger
	loggerCfg *config.Logging
}

func Create(cfg *config.Logging) *AppLogger {
	return &AppLogger{loggerCfg: cfg, logger: log.New()}
}

func (l *AppLogger) Info(message string) {
	l.logger.WithFields(log.Fields{"location": l.getLocation()}).Info(message)
}

func (l *AppLogger) Warn(message string) {
	l.logger.WithFields(log.Fields{"location": l.getLocation()}).Warn(message)
}

func (l *AppLogger) Error(message string) {
	l.logger.WithFields(log.Fields{"location": l.getLocation()}).Error(message)
}

func (l *AppLogger) Debug(message string) {
	l.logger.WithFields(log.Fields{"location": l.getLocation()}).Debug(message)
}

func (l *AppLogger) Trace(message string) {
	l.logger.WithFields(log.Fields{"location": l.getLocation()}).Trace(message)
}

func (l *AppLogger) Init() {
	l.logger.Out = io.Discard
	for _, a := range l.loggerCfg.Appenders {
		if !a.Enabled {
			continue
		}

		level := l.getLevel(*l.loggerCfg.Level)
		if a.Level != nil {
			level = min(level, l.getLevel(*a.Level))
		}

		switch a.Type {
		case config.RollingFile:
			hook, _ := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
				Filename:   string(a.Destination.File),
				MaxSize:    a.Destination.MaxSize,
				MaxBackups: a.Destination.MaxBackups,
				MaxAge:     a.Destination.MaxAge,
				Level:      level,
				Formatter: &log.TextFormatter{
					FullTimestamp:   true,
					TimestampFormat: timestampFormat,
				},
			})
			l.logger.AddHook(hook)

		case config.Console:
			l.logger.SetLevel(level)
			l.logger.SetOutput(colorable.NewColorableStdout())
			l.logger.SetFormatter(&log.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: timestampFormat,
			})
		}
	}
}

func (l *AppLogger) getLocation() string {
	// runtime.Caller ascends two stack frames to get to the appropriate location
	// and return valid line from the code
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "debug logger now its broken"
		line = 0
		return fmt.Sprintf("%s:%d :", filepath.Base(file), line)
	}
	return fmt.Sprintf("%s:%d :", filepath.Base(file), line)
}

func (l *AppLogger) getLevel(level string) log.Level {
	lev, ok := logLevels[level]
	if ok {
		return lev
	}
	return defaultLogLevel
}

func min(x, y log.Level) log.Level {
	if x < y {
		return x
	}

	return y
}
