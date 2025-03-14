package tools

import (
	"context"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/zhangga/aino/internal/tools/larkcli"
	"github.com/zhangga/aino/pkg/logger"
)

var _ tool.InvokableTool = (*ToolLarkGetMsg)(nil)

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
	err := json.Unmarshal([]byte(argumentsInJSON), p)
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
