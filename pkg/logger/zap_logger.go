package logger

import (
	"os"
	config "scs-user/config"
	"sync" // Import the sync package

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger methods interface
type Logger interface {
	InitLogger(cfg *config.Config) // Fixed to match implementation
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	IsInitialized() bool // Add method to check initialization status
}

// ApiLogger
type ApiLogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
	initialized bool
}

var (
	once sync.Once
	log  *ApiLogger
)

// GetLogger returns the singleton instance of the logger.
// It initializes the logger only on the first call.
func GetLogger() *ApiLogger {
	once.Do(func() {
		log = &ApiLogger{} // The config will be set in InitLogger
	})
	return log
}

// InitLogger initializes the logger with the provided configuration.
// It is intended to be called once, typically in your main function.
func (l *ApiLogger) InitLogger(cfg *config.Config) {
	l.cfg = cfg
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Mode == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
	l.initialized = true
}

// IsInitialized returns true if the logger has been initialized
func (l *ApiLogger) IsInitialized() bool {
	return l.initialized && l.sugarLogger != nil
}

// ... (rest of your logging methods like Debug, Info, etc.)

// getLoggerLevel and other helper methods remain the same
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *ApiLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}
func (l *ApiLogger) Debug(args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Debug(args...)
}

func (l *ApiLogger) Debugf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Debugf(template, args...)
}

func (l *ApiLogger) Info(args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Info(args...)
}

func (l *ApiLogger) Infof(template string, args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Infof(template, args...)
}

func (l *ApiLogger) Warn(args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Warn(args...)
}

func (l *ApiLogger) Warnf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Warnf(template, args...)
}

func (l *ApiLogger) Error(args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Error(args...)
}

func (l *ApiLogger) Errorf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.Errorf(template, args...)
}

func (l *ApiLogger) DPanic(args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.DPanic(args...)
}

func (l *ApiLogger) DPanicf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		return
	}
	l.sugarLogger.DPanicf(template, args...)
}

func (l *ApiLogger) Panic(args ...interface{}) {
	if !l.IsInitialized() {
		panic("Logger not initialized")
	}
	l.sugarLogger.Panic(args...)
}

func (l *ApiLogger) Panicf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		panic("Logger not initialized")
	}
	l.sugarLogger.Panicf(template, args...)
}

func (l *ApiLogger) Fatal(args ...interface{}) {
	if !l.IsInitialized() {
		os.Exit(1)
	}
	l.sugarLogger.Fatal(args...)
}

func (l *ApiLogger) Fatalf(template string, args ...interface{}) {
	if !l.IsInitialized() {
		os.Exit(1)
	}
	l.sugarLogger.Fatalf(template, args...)
}
