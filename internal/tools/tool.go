package tools

import (
	"context"
	"github.com/cloudwego/eino/components/tool"
	"github.com/zhangga/aino/internal/conf"
)

var (
	appCtx    context.Context
	appConfig *conf.Config
)

func InitTools(ctx context.Context, config *conf.Config) {
	appCtx = ctx
	appConfig = config
}

type Creator func() tool.BaseTool
