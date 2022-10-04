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

//func (l *AppLogger) Log(level log.Level, message string) {
//l.logger.WithFields(log.Fields{ "location": l.getLocation(),
//}).Info(message)
//}

func (l *AppLogger) Init() {
	l.logger.Out = io.Discard
	for _, a := range l.loggerCfg.Appenders {
		if !a.Enabled {
			continue
		}

		level := GetLevel(l.loggerCfg.Level)
		if a.Level != nil {
			level = Min(level, GetLevel(a.Level))
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

// Defining logger from logrus here
//var logger = log.New()
//var LoggerConfig = config.GlobalConfig{}

func GetLevel(level *string) log.Level {
	switch *level {
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.TraceLevel
	}
}

func Min(x, y log.Level) log.Level {
	if x < y {
		return x
	}

	return y
}

func GetAppenderIndex(appenderType config.AppenderType, appenders []config.Appender) int {
	for i, v := range appenders {
		if v.Type == appenderType {
			return i
		}
	}

	return -1
}

// InitLoggerFromByteSlice initialize logger by data byte slice.
/*func InitLoggerFromByteSlice(configData []byte) utils.ShutdownFunc {
	config, err := setup.ReadConfigFromByteSlice[config.GlobalConfig](configData)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	LoggerConfig = *config

	fmt.Println(LoggerConfig)

	return setupLogger()
}*/

// InitLoggerFromFile initialize logger by config file, stored in configPath file.
/*func InitLoggerFromFile(configPath string) utils.ShutdownFunc {
	absPath, _ := filepath.Abs(configPath)

	fileData, err := ioutil.ReadFile(absPath)
	if err != nil {
		return setupLoggerWhenError(err)
	}

	return InitLoggerFromByteSlice(fileData)
}*/

// setupLoggerWhenError creates simple logger to stdout.
/*func setupLoggerWhenError(err error) utils.ShutdownFunc {
	fmt.Println("Error reading loggerconfig:", err)
	fmt.Println("Setting logger output to console, level to trace")

	logger.Formatter = &log.TextFormatter{}
	logger.Out = os.Stdout
	logger.Level = log.TraceLevel

	ErrorLog(err.Error())

	return utils.NewShutdowner().Shutdown
}*/

// setupLogger create appeanders. Execute only if loggerConfig is setup.
/*func setupLogger() utils.ShutdownFunc {
	shutdowner := utils.NewShutdowner()

	logger.Out = io.Discard
	for _, a := range LoggerConfig.Logging.Appenders {
		if !a.Enabled {
			continue
		}

		level := GetLevel(LoggerConfig.Logging.Level)
		if a.Level != nil {
			level = Min(level, GetLevel(a.Level))
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
					TimestampFormat: LoggerTimestampFormat,
				},
			})
			logger.AddHook(hook)

		case config.Console:
			logger.SetLevel(level)
			logger.SetOutput(colorable.NewColorableStdout())
			logger.SetFormatter(&log.TextFormatter{
				ForceColors:     true,
				FullTimestamp:   true,
				TimestampFormat: LoggerTimestampFormat,
			})
		}
	}

	return shutdowner.Shutdown
}*/

/*func InfoLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Info(message)
}

func WarnLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Warn(message)
}

func FatalLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Fatal(message)
}

func PanicLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Panic(message)
}

func ErrorLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Error(message)
}

func TraceLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Trace(message)
}

func DebugLog(message string) {
	logger.WithFields(log.Fields{
		"location": location(),
	}).Debug(message)
}

// Location finds location of the log call and returns it as a string.
func location() string {
	// runtime.Caller ascends two stack frames to get to the appropriate location
	// and return valid line from the code
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "debug logger now its broken"
		line = 0
		return fmt.Sprintf("%s:%d :", filepath.Base(file), line)
	}
	return fmt.Sprintf("%s:%d :", filepath.Base(file), line)
}*/
