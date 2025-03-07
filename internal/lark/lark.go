package lark

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func RunService(ctx context.Context, appId, appSecret string) {
	// 注册事件回调，OnP2MessageReceiveV1 为接收消息 v2.0；OnCustomizedEvent 内的 message 为接收消息 v1.0。
	eventHandler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			// 获取父消息ID
			// parentMsgId := event.Message.ParentId
			// if parentMsgId != "" {
			// 	// 调用获取消息详情API
			// 	resp, err := larkim.NewService(cli.Client).Messages.Get(context.Background(),
			// 		larkim.NewGetMessageReqBuilder()
			// 			.MessageId(parentMsgId)
			// 			.Build())

			// 	if err != nil {
			// 		fmt.Printf("获取父消息失败: %v\n", err)
			// 	} else {
			// 		fmt.Printf("被回复的原始消息内容: %s\n",
			// 			larkcore.Prettify(resp.Items[0].Body.Content))
			// 	}
			// }

			fmt.Printf("[ OnP2MessageReceiveV1 access ], data: %s\n", larkcore.Prettify(event))
			return nil
		}).
		OnCustomizedEvent("这里填入你要自定义订阅的 event 的 key，例如 out_approval", func(ctx context.Context, event *larkevent.EventReq) error {
			fmt.Printf("[ OnCustomizedEvent access ], type: message, data: %s\n", string(event.Body))
			return nil
		})

	// 创建Client
	cli := larkws.NewClient(appId, appSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
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
