package main

import (
	"context"
	"github.com/zhangga/aino/internal/eino"
	"github.com/zhangga/aino/internal/rolelist"
	"github.com/zhangga/aino/internal/service"
	"github.com/zhangga/aino/internal/tools"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/zhangga/aino/pkg/logger"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/internal/larksrv"
	"github.com/zhangga/aino/pkg/config"
)

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "run the aino service",
	Run:   run,
}

var (
	ConfigPath string
	Config     conf.Config
)

func init() {
	cmdRun.Flags().StringVarP(&ConfigPath, "config", "c", "configs/config.yaml", "config file path")

	// 需要绑定命令行参数的时候可以在这设置，如:
	Config.LarkConfig = &conf.LarkConfig{}
	cmdRun.Flags().StringVar(&Config.LarkConfig.AppID, "lark.app_id", "", "Lark AppId. ENV: LARK_APP_ID")
	cmdRun.Flags().StringVar(&Config.LarkConfig.AppSecret, "lark.app_secret", "", "Lark AppSecret. ENV: LARK_APP_SECRET")

	Config.LLMConfig = &conf.LLMConfig{}
	cmdRun.Flags().StringVar(&Config.LLMConfig.Model, "llm.model", "", "LLM Model. ENV: LLM_MODEL")
	cmdRun.Flags().StringVar(&Config.LLMConfig.ApiKey, "llm.api_key", "", "LLM ApiKey. ENV: LLM_API_KEY")
	cmdRun.Flags().StringVar(&Config.LLMConfig.BaseURL, "llm.base_url", "", "LLM BaseUrl. ENV: LLM_BASE_URL")
}

func init() {
	// 启用环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 绑定命令行参数
	if err := viper.BindPFlags(cmdRun.Flags()); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// 初始化日志
	logger.InitLogger()
	defer logger.Sync()

	logger.Info("Aino service starting...")

	// 指定了配置文件，但是文件不存在
	if cmd.Flags().Changed("config") {
		// 检查文件是否存在
		if _, err := os.Stat(ConfigPath); err != nil {
			logger.Fatal("config file not exist",
				zap.String("path", ConfigPath),
				zap.Error(err))
		}
	}

	// 读取配置
	if err := config.LoadConfig(&Config, ConfigPath); err != nil {
		logger.Fatalf("load config file: %s failed, err=%v", ConfigPath, err)
	}

	// 控制组件的启动和关闭
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rolelist.InitRoleList(ctx)
	tools.InitTools(ctx, &Config)

	// 启动处理服务
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := service.RunService(ctx, Config.ServiceConfig.HttpPort); err != nil {
			logger.Fatal("run service failed", zap.Error(err))
		}
	}()

	// 启动LarkService
	wg.Add(1)
	go func() {
		defer wg.Done()
		larksrv.RunService(ctx, Config.LarkConfig.AppID, Config.LarkConfig.AppSecret)
	}()

	// 启动langchain
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	langchain.Run(ctx, Config.LLMConfig)
	//}()

	// 初始化eino agent
	if err := eino.InitAgent(ctx, Config.LLMConfig); err != nil {
		logger.Fatal("init eino agent failed, err: %v", err)
	}

	// 监听系统信号
	listenForOSSignal(ctx, cancel, &wg)

	// 等待结束
	shutdownGracefully(ctx, &wg)
}

func listenForOSSignal(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		select {
		case s := <-chSignal:
			logger.Infof("os.Signal received: %s\n", s.String())
		case <-ctx.Done():
			return
		}

		// shutdown
		signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		cancel()
	}()
}

func shutdownGracefully(_ context.Context, wg *sync.WaitGroup) {
	wg.Wait()
}
