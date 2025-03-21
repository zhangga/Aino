package langchain

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/pkg/logger"
)

func Run(ctx context.Context, llmConfig *conf.LLMConfig) {
	// Initialize the OpenAI client with Deepseek model
	llm, err := openai.New(
		openai.WithModel(llmConfig.Model),
		openai.WithToken(llmConfig.ApiKey),
		openai.WithBaseURL(llmConfig.BaseURL),
	)
	if err != nil {
		logger.Fatal(err)
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "你是人工智能助手"),
		llms.TextParts(llms.ChatMessageTypeHuman, "你好"),
	}
	if _, err := llm.GenerateContent(ctx, content,
		llms.WithMaxTokens(1024),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		})); err != nil {
		logger.Errorf("langchain.Generate failed, err=%v", err)
	}

}
