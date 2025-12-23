package aino

import (
	"context"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CmdRun = &cobra.Command{
	Use:   "aino",
	Short: "run the aino service",
	Run:   run,
}

var (
	ConfigPath string
)

func init() {
	CmdRun.Flags().StringVarP(&ConfigPath, "config", "c", "configs/config.yaml", "config file path")
}

func init() {
	// 启用环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 绑定命令行参数
	if err := viper.BindPFlags(CmdRun.Flags()); err != nil {
		panic(err)
	}
}

func run(cmd *cobra.Command, args []string) {

}

func listenForOSSignal(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	wg.Add(1)
}

func shutdownGracefully(_ context.Context, wg *sync.WaitGroup) {
	wg.Wait()
}
