package zlog

import (
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Logger = (*zap.SugaredLogger)(nil)

var (
	seedLogger    *zap.Logger
	defaultLogger Logger
	defaultInit   sync.Once
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Sync() error
}

func InitDefaultLogger(config Config) {
	defaultInit.Do(func() {
		writeSyncer := getLogWriter(config)
		encoder := getEncoder()
		level, err := zapcore.ParseLevel(config.Level)
		if err != nil {
			level = zapcore.DebugLevel
		}
		core := zapcore.NewCore(encoder, writeSyncer, level)

		seedLogger = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zap.PanicLevel),
		)
		defaultLogger = seedLogger.WithOptions(
			zap.AddCallerSkip(1),
		).Sugar()
	})
}

// NewLogger creates a new Logger based on the provided configuration.
func NewLogger(config Config) Logger {
	writeSyncer := getLogWriter(config)
	encoder := getEncoder()
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.PanicLevel),
	)
	return logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func Sync() {
	if err := defaultLogger.Sync(); err != nil {
		panic(err)
	}
}

// WithOptions creates a new Logger based on the seedLogger with the provided options.
func WithOptions(opts ...zap.Option) Logger {
	return seedLogger.WithOptions(opts...).Sugar()
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(config Config) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    100, //MB
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   false,
	}
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}
