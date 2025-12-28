package tools

import (
	"github.com/cloudwego/eino/components/tool"
)

type Tool interface {
	ToEinoTool() (tool.BaseTool, error)
}
