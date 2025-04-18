package tools

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/zhangga/aino/internal/tools/larkcli"
	"github.com/zhangga/aino/pkg/logger"
)

var (
	_ tool.InvokableTool = (*ToolLarkGetMsg)(nil)
	_ tool.InvokableTool = (*ToolLarkSendMsg)(nil)
)

type ToolLarkGetMsg struct {
	larkCli *larkcli.Client
}

type ToolLarkGetMsgParam struct {
	MessageId string `json:"message_id"`
}

func NewToolLarkGetMsg() tool.BaseTool {
	return &ToolLarkGetMsg{
		larkCli: larkcli.NewClient(appCtx, appConfig.LarkConfig.AppID, appConfig.LarkConfig.AppSecret),
	}
}

func (t *ToolLarkGetMsg) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_lark_message",
		Desc: "Get lark message",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"message_id": {
				Type:     "string",
				Desc:     "The lark message id",
				Required: true,
			},
		}),
	}, nil
}

func (t *ToolLarkGetMsg) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &ToolLarkGetMsgParam{}
	err := sonic.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}

	message, err := t.larkCli.GetMessage(p.MessageId)
	if err != nil {
		return "", err
	}
	json, err := sonic.Marshal(message)
	if err != nil {
		return "", err
	}
	logger.Debugf("get lark message: %s", json)
	return string(json), nil
}

type ToolLarkSendMsg struct {
	larkCli *larkcli.Client
}

type ToolLarkSendMsgParam struct {
	ChatId  string `json:"chat_id"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

func NewToolLarkSendMsg() tool.BaseTool {
	return &ToolLarkSendMsg{
		larkCli: larkcli.NewClient(appCtx, appConfig.LarkConfig.AppID, appConfig.LarkConfig.AppSecret),
	}
}

func (t *ToolLarkSendMsg) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "send_lark_message",
		Desc: "Send lark message",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"chat_id": {
				Type:     "string",
				Desc:     "The lark chat id",
				Required: true,
			},
			"uuid": {
				Type:     "string",
				Desc:     "The lark message uuid",
				Required: true,
			},
			"content": {
				Type:     "string",
				Desc:     "The lark message content",
				Required: true,
			},
		}),
	}, nil
}

func (t *ToolLarkSendMsg) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &ToolLarkSendMsgParam{}
	err := sonic.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}
	err = t.larkCli.SendMessage(p.ChatId, "text", p.Content, p.UUID)
	if err != nil {
		return "", err
	}
	return "success", nil
}

func (t *ToolLarkSendMsg) SendMessage(ctx context.Context, chatId, content string) error {
	return t.larkCli.SendMessage(chatId, "text", content, uuid.NewString())
}
