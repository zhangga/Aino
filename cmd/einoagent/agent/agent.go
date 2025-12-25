package agent

import (
	"context"
	"io"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	"github.com/zhangga/aino/conf"
	logger "github.com/zhangga/aino/pkg/zlog"
)

var (
	once sync.Once
	log  logger.Logger
)

func Init() error {
	var err error
	once.Do(func() {
		cbConfig := &LogCallbackConfig{
			Detail: true,
			Debug:  conf.GlobalConfig.ServiceConf.Debug,
		}
		logger.WithOptions()
	})
	return err
}

type LogCallbackConfig struct {
	Detail bool
	Debug  bool
}

func LogCallback(config *LogCallbackConfig) callbacks.Handler {
	if config == nil {
		config = &LogCallbackConfig{
			Detail: true,
			Debug:  false,
		}
	}
	builder := callbacks.NewHandlerBuilder()
	builder.OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {

	})
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		log.Infof("[view]: end [%s:%s:%s]", info.Component, info.Type, info.Name)
		return ctx
	})
	return builder.Build()
}
