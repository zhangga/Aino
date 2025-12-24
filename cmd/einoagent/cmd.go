package einoagent

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/devops"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/cobra"
	"github.com/zhangga/aino/cmd/einoagent/agent"
	"github.com/zhangga/aino/cmd/einoagent/task"
	"github.com/zhangga/aino/conf"
	"github.com/zhangga/aino/pkg/utils"
	logger "github.com/zhangga/aino/pkg/zlog"
)

var CmdRun = &cobra.Command{
	Use:   "einoagent",
	Short: "run the einoagent service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logger.InitLogConfig(logger.Config{FilePath: "logs/einoagent.log", Level: "debug"})
	defer logger.Sync()

	logger.Info("starting einoagent service...")
	ctx := context.Background()

	// Enable Eino Debug Mode if configured
	if conf.GlobalConfig.ServiceConf.EinoDebug {
		logger.Info("Eino Debug Mode is enabled.")
		if err := devops.Init(ctx); err != nil {
			logger.Errorf("Failed to initialize devops: %s", err)
		}
	}

	// 创建 Hertz 服务器并运行
	h := server.Default(server.WithHostPorts(fmt.Sprintf(":%d", conf.GlobalConfig.ServiceConf.HttpPort)))
	h.Use(LogMiddleware())

	//TODO APMPLUS
	if len(conf.GlobalConfig.ServiceConf.APMPlusAppKey) > 0 {
		logger.Info("APMPlus is enabled.")
		// apmplus.InitAPMPlus(conf.GlobalConfig.ServiceConf.APMPlusAppKey, h)
	}

	// 注册 task 路由
	taskGroup := h.Group("/task")
	if err := task.BindRoutes(taskGroup); err != nil {
		logger.Fatalf("failed to bind task routes: %v", err)
	}
	// 注册 agent 路由
	agentGroup := h.Group("/agent")
	if err := agent.BindRoutes(agentGroup); err != nil {
		logger.Fatalf("failed to bind agent routes: %v", err)
	}
	// Redirect root path to /agent
	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Redirect(302, []byte("/agent"))
	})

	// 启动服务器
	h.Spin()
}
func LogMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := utils.BytesToString(c.Request.URI().Path())
		method := utils.BytesToString(c.Request.Header.Method())

		// 处理请求
		c.Next(ctx)

		// 记录请求信息
		latency := time.Since(start)
		status := c.Response.StatusCode()
		logger.Infof("[HTTP] %s %s %d %s", method, path, status, latency)
	}
}
