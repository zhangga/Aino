package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

var _ ILogger = (*zap.SugaredLogger)(nil)

var Default ILogger

type ILogger interface {
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

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.PanicLevel))
	Default = logger.Sugar()
}

func Sync() {
	if err := Default.Sync(); err != nil {
		panic(err)
	}
}

func Debug(args ...interface{}) {
	Default.Debug(args...)
}

func Info(args ...interface{}) {
	Default.Info(args...)
}

func Warn(args ...interface{}) {
	Default.Warn(args...)
}

func Error(args ...interface{}) {
	Default.Error(args...)
}

func Panic(args ...interface{}) {
	Default.Panic(args...)
}

func Fatal(args ...interface{}) {
	Default.Fatal(args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Default.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Default.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Default.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Default.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	Default.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	Default.Fatalw(msg, keysAndValues...)
}

func Debugf(template string, args ...interface{}) {
	Default.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Default.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Default.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Default.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	Default.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Default.Fatalf(template, args...)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    100, //MB
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	ws := io.MultiWriter(lumberJackLogger, os.Stdout)
	return zapcore.AddSync(ws)
}
