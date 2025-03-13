package eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/pkg/logger"
)

func Run(ctx context.Context, llmConfig *conf.LLMConfig) (string, error) {
	config := &openai.ChatModelConfig{
		Model:   llmConfig.Model,
		APIKey:  llmConfig.ApiKey,
		BaseURL: llmConfig.BaseURL,
	}
	model, _ := openai.NewChatModel(ctx, config) // create an invokable LLM instance

	message, err := model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("you are a helpful assistant."),
		schema.UserMessage("what does the future AI App look like?"),
	})
	if err != nil {
		logger.Errorf("eino.Generate failed, err=%v", err)
		return "", err
	}

	logger.Info(message.Content)
	return "", nil
}
