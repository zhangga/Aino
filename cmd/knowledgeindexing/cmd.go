package knowledgeindexing

import (
	"github.com/spf13/cobra"
	logger "github.com/zhangga/aino/pkg/zlog"
)

var CmdRun = &cobra.Command{
	Use:   "knowledge",
	Short: "run the knowledgeindexing service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logger.InitLogConfig(logger.Config{"logs/knowledgeindexing.log", "debug"})
	defer logger.Sync()

	logger.Info("starting knowledgeindexing service...")
}
