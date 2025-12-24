package einoagent

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// newLambda1 component initialization function of node 'Lambda2' in graph 'EinoAgent'
func newLambda1(ctx context.Context) (lba *compose.Lambda, err error) {
	// TODO Modify component configuration here.
	config := &react.AgentConfig{}
	chatModelIns11, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	config.Model = chatModelIns11
	toolIns21, err := newTool(ctx)
	if err != nil {
		return nil, err
	}
	toolIns22, err := newTool1(ctx)
	if err != nil {
		return nil, err
	}
	toolIns23, err := newTool2(ctx)
	if err != nil {
		return nil, err
	}
	toolIns24, err := newTool3(ctx)
	if err != nil {
		return nil, err
	}
	toolIns25, err := newTool4(ctx)
	if err != nil {
		return nil, err
	}
	config.ToolsConfig.Tools = []tool.BaseTool{toolIns21, toolIns22, toolIns23, toolIns24, toolIns25}
	ins, err := react.NewAgent(ctx, config)
	if err != nil {
		return nil, err
	}
	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)
	if err != nil {
		return nil, err
	}
	return lba, nil
}
