package logger_test

import (
	"github.com/zhangga/aino/pkg/logger"
	"go.uber.org/zap"
	"testing"
)

func TestLoggerSkip(t *testing.T) {
	log := logger.WithOptions(zap.AddCallerSkip(2))
	//defer func() {
	//	_ = log.Sync()
	//}()

	type User struct {
		Name string
		Age  int
	}

	skipLog1 := func(log logger.ILogger, v interface{}) {
		log.Debugf("this is SkipLog: %v", v)
	}

	skipLog2 := func(log logger.ILogger, v interface{}) {
		skipLog1(log, v)
	}

	skipLog2(log, User{Name: "John", Age: 42})
}
