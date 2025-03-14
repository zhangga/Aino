package handler

import (
	"context"
	_ "embed"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/zhangga/aino/internal/eino"
	"github.com/zhangga/aino/internal/tools"
	"github.com/zhangga/aino/pkg/logger"
)

//go:embed larkprompt.txt
var larkPrompt string

var larkTools = []tools.Creator{
	tools.NewToolLarkGetMsg,
}

var _ Task = (*LarkTask)(nil)

type LarkTask struct {
	id  uint64
	Msg *larkim.P2MessageReceiveV1
}

func NewLarkTask(event *larkim.P2MessageReceiveV1) *LarkTask {
	task := &LarkTask{
		id:  NextID(),
		Msg: event,
	}
	return task
}

func (t *LarkTask) Id() uint64 {
	return t.id
}

func (t *LarkTask) Type() TaskType {
	return TaskTypeLark
}

func (t *LarkTask) Run(ctx context.Context) {
	ragent, err := eino.NewAgentByType(ctx, larkPrompt, larkTools...)
	if err != nil {
		panic(err)
	}

	message, err := sonic.Marshal(t.Msg.Event)
	if err != nil {
		panic(err)
	}
	result, err := ragent.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: string(message),
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&eino.LoggerCallback{})))
	if err != nil {
		panic(err)
	}

	logger.Infof("Lark taskId=%d, result=%v", t.id, result)
}
