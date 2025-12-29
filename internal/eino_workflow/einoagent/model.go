package einoagent

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
	"github.com/zhangga/aino/conf"
)

func newToolCallingChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ark.ChatModelConfig{
		BaseURL: conf.GlobalConfig.LLMConf.BaseURL,
		Model:   conf.GlobalConfig.LLMConf.Model,
		APIKey:  conf.GlobalConfig.LLMConf.APIKey,
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
