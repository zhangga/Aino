package eino

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/internal/tools"
)

var (
	llmConfig *conf.LLMConfig
)

func InitAgent(_ context.Context, llmConf *conf.LLMConfig) error {
	llmConfig = llmConf
	return nil
}

func NewAgentByType(ctx context.Context, persona string, creators ...tools.Creator) (*react.Agent, error) {
	var ts []tool.BaseTool
	for _, c := range creators {
		ts = append(ts, c())
	}
	return NewAgent(ctx, persona, ts...)
}

func NewAgent(ctx context.Context, persona string, tools ...tool.BaseTool) (*react.Agent, error) {
	// create an invokable LLM instance
	model, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   llmConfig.Model,
		APIKey:  llmConfig.ApiKey,
		BaseURL: llmConfig.BaseURL,
	})
	if err != nil {
		return nil, err
	}

	config := &react.AgentConfig{
		Model:           model,
		ToolsConfig:     compose.ToolsNodeConfig{Tools: tools},
		MessageModifier: react.NewPersonaModifier(persona),
		// StreamToolCallChecker: toolCallChecker, // uncomment it to replace the default tool call checker with custom one
	}
	return react.NewAgent(ctx, config)
}
