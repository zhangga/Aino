package agent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/callbacks/apmplus"
	"github.com/cloudwego/eino-ext/callbacks/langfuse"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/zhangga/aino/conf"
	"github.com/zhangga/aino/internal/eino_workflow/einoagent"
	"github.com/zhangga/aino/pkg/mempkg"
	"github.com/zhangga/aino/pkg/utils"
	logger "github.com/zhangga/aino/pkg/zlog"
)

const (
	APMPlusHost  = "apmplus-%s.volces.com:4317"
	langfuseHost = "https://cloud.langfuse.com"
)

var (
	once        sync.Once
	agentLogger logger.Logger
	cbHandler   callbacks.Handler
	memory      = mempkg.GetDefaultMemory()
)

func Init() error {
	var err error
	once.Do(func() {
		agentLogger = logger.NewLogger(logger.Config{
			Level:    "debug",
			FilePath: "logs/einoagent_detail.log",
		})

		cbConfig := &LogCallbackConfig{
			Detail: true,
			Debug:  conf.GlobalConfig.ServiceConf.Debug,
		}
		cbHandler = LogCallback(cbConfig)

		// init global callback, for trace and metrics
		callbackHandlers := make([]callbacks.Handler, 0)
		if len(conf.GlobalConfig.ServiceConf.APMPlusAppKey) > 0 {
			agentLogger.Info("[eino agent]: use apmplus as callback, watch at: https://console.volcengine.com/apmplus-server")
			region := conf.GlobalConfig.ServiceConf.APMPlusRegion
			cbh, _, err := apmplus.NewApmplusHandler(&apmplus.Config{
				Host:        fmt.Sprintf(APMPlusHost, region),
				AppKey:      conf.GlobalConfig.ServiceConf.APMPlusAppKey,
				ServiceName: "eino-assistant",
				Release:     "release/v0.0.1",
			})
			if err != nil {
				agentLogger.Fatalf("init apmplus callback handler failed: %v", err)
			}
			callbackHandlers = append(callbackHandlers, cbh)
		}

		if len(conf.GlobalConfig.ServiceConf.LangfusePublicKey) > 0 && len(conf.GlobalConfig.ServiceConf.LangfuseSecretKey) > 0 {
			agentLogger.Infof("[eino agent]: use langfuse as callback, watch at: https://cloud.langfuse.com")
			cbh, _ := langfuse.NewLangfuseHandler(&langfuse.Config{
				Host:      langfuseHost,
				PublicKey: conf.GlobalConfig.ServiceConf.LangfusePublicKey,
				SecretKey: conf.GlobalConfig.ServiceConf.LangfuseSecretKey,
				Name:      "Eino Assistant",
				Public:    true,
				Release:   "release/v0.0.1",
				UserID:    "eino_god",
				Tags:      []string{"eino", "assistant"},
			})
			callbackHandlers = append(callbackHandlers, cbh)
		}
		if len(callbackHandlers) > 0 {
			callbacks.AppendGlobalHandlers(callbackHandlers...)
		}
	})
	return err
}

func RunAgent(ctx context.Context, id, msg string) (*schema.StreamReader[*schema.Message], error) {
	runner, err := einoagent.BuildEinoAgent(ctx)
	if err != nil {
		return nil, fmt.Errorf("build eino agent failed: %w", err)
	}

	conversation := memory.GetConversation(id, true)

	userMessage := &einoagent.UserMessage{
		Id:      id,
		Query:   msg,
		History: conversation.GetMessages(),
	}
	if len(conf.GlobalConfig.ServiceConf.APMPlusAppKey) > 0 {
		ctx = apmplus.SetSession(ctx, apmplus.WithSessionID(id), apmplus.WithUserID("eino-assistant-user"))
	}
	sr, err := runner.Stream(ctx, userMessage, compose.WithCallbacks(cbHandler))
	if err != nil {
		return nil, fmt.Errorf("run eino agent failed: %w", err)
	}

	srs := sr.Copy(2)

	go func() {
		// for save to memory
		fullMsgs := make([]*schema.Message, 0)

		defer func() {
			// close stream if you used it
			srs[1].Close()

			// add user input to history
			conversation.Append(schema.UserMessage(msg))

			fullMsg, err := schema.ConcatMessages(fullMsgs)
			if err != nil {
				agentLogger.Errorf("error concatenating messages: %s", err.Error())
			}
			// add agent response to history
			conversation.Append(fullMsg)
		}()

	outer:
		for {
			select {
			case <-ctx.Done():
				agentLogger.Infof("context done: %v", ctx.Err())
				return
			default:
				chunk, err := srs[1].Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break outer
					}
				}

				fullMsgs = append(fullMsgs, chunk)
			}
		}
	}()

	return srs[0], nil
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
		agentLogger.Infof("[view]: start [%s:%s:%s]", info.Component, info.Type, info.Name)
		if config.Detail {
			var b []byte
			if config.Debug {
				b, _ = sonic.MarshalIndent(input, "", "  ")
			} else {
				b, _ = sonic.Marshal(input)
			}
			agentLogger.Infof("[view]: input: %s", utils.BytesToString(b))
		}
		return ctx
	})
	builder.OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		agentLogger.Infof("[view]: end [%s:%s:%s]", info.Component, info.Type, info.Name)
		return ctx
	})
	return builder.Build()
}
