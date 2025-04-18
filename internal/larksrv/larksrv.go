package larksrv

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	"github.com/zhangga/aino/internal/service"
	"github.com/zhangga/aino/pkg/logger"
)

// RunService 启动服务监听lark消息
func RunService(ctx context.Context, appId, appSecret string) {
	// 注册事件回调，OnP2MessageReceiveV1 为接收消息 v2.0；OnCustomizedEvent 内的 message 为接收消息 v1.0。
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			logger.Debugf("[ Lark.OnP2MessageReceiveV1 access ], data: %s\n", larkcore.Prettify(event))
			// 创建处理lark消息的任务
			task := service.NewLarkTask(event)
			if err := service.AddTask(task); err != nil {
				logger.Errorf("[ Lark.OnP2MessageReceiveV1 access ], add task err: %v\n", err)
				return err
			}
			return nil
		}).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			logger.Debugf("[ Lark.OnP2MessageReadV1 access ], data: %s\n", larkcore.Prettify(event))
			return nil
		}).
		OnCustomizedEvent("这里填入你要自定义订阅的 event 的 key，例如 out_approval", func(ctx context.Context, event *larkevent.EventReq) error {
			logger.Debugf("[ Lark.OnCustomizedEvent access ], type: message, data: %s\n", string(event.Body))
			return nil
		})

	// 创建Client
	cli := larkws.NewClient(appId, appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
		larkws.WithLogger(&LarkLogger{Logger: logger.Default}),
	)

	// 启动客户端。目前lark client里不响应ctx的cancel，自己来控制了
	ch := make(chan struct{})
	go func() {
		defer close(ch)

		if err := cli.Start(ctx); err != nil {
			panic(err)
		}
	}()

	select {
	case <-ctx.Done():
	case <-ch:
	}
}

type LarkLogger struct {
	Logger logger.ILogger
}

func (z *LarkLogger) Debug(ctx context.Context, args ...interface{}) {
	z.Logger.Debug(fmt.Sprint(args...))
}

func (z *LarkLogger) Info(ctx context.Context, args ...interface{}) {
	z.Logger.Info(fmt.Sprint(args...))
}

func (z *LarkLogger) Warn(ctx context.Context, args ...interface{}) {
	z.Logger.Warn(fmt.Sprint(args...))
}

func (z *LarkLogger) Error(ctx context.Context, args ...interface{}) {
	z.Logger.Error(fmt.Sprint(args...))
}
