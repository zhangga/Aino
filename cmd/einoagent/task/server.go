package task

import (
	"context"
	"embed"

	"github.com/cloudwego/hertz/pkg/route"
)

//go:embed web
var webContent embed.FS

// BindRoutes 注册路由
func BindRoutes(r *route.RouterGroup) error {
	ctx := context.Background()
	_ = ctx

	//taskTool, err := task.NewTaskToolImpl(ctx, &task.TaskToolConfig{
	//	Storage: task.GetDefaultStorage(),
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}
