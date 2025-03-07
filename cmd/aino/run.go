package main

import (
	"context"
	"github.com/zhangga/aino/internal/eino"
	"github.com/zhangga/aino/internal/langchain"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/internal/lark"
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
	// 指定了配置文件，但是文件不存在
	if cmd.Flags().Changed("config") {
		// 检查文件是否存在
		if _, err := os.Stat(ConfigPath); err != nil {
			log.Fatalf("config file: %s, no exist, err=%v", ConfigPath, err)
		}
	}

	// 读取配置
	if err := config.LoadConfig(&Config, ConfigPath); err != nil {
		log.Fatalf("load config file: %s failed, err=%v", ConfigPath, err)
	}

	// 控制组件的启动和关闭
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动LarkService，目前lark里面没法友好结束，先不用wg控制
	wg.Add(1)
	go func() {
		defer wg.Done()

		lark.RunService(ctx, Config.LarkConfig.AppID, Config.LarkConfig.AppSecret)
	}()

	// 启动langchain
	wg.Add(1)
	go func() {
		defer wg.Done()

		langchain.Run(ctx, Config.LLMConfig)
	}()

	// 启动eino
	wg.Add(1)
	go func() {
		defer wg.Done()

		eino.Run(ctx, Config.LLMConfig)
	}()

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
			log.Printf("os.Signal received: %s\n", s.String())
		case <-ctx.Done():
			return
		}

		// shutdown
		signal.Reset(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		cancel()
	}()
}

func shutdownGracefully(ctx context.Context, wg *sync.WaitGroup) {
	wg.Wait()
}
