package einoagent

import "github.com/spf13/cobra"

var CmdRun = &cobra.Command{
	Use:   "einoagent",
	Short: "run the einoagent service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {

}
