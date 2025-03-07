package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the project.
	Version = "0.0.1"
	// BuildTime is the time when the project was built.
	BuildTime = "unknown"
	// CommitHash is the commit hash of the project.
	CommitHash = "unknown"
	// Description is the description of the project.
	Description = "Aino is a tool for managing version."
)

// Print 打印版本信息
func Print() string {
	kv := [][2]string{
		{"Go", runtime.Version()},
		{"Version", Version},
		{"BuildTime", BuildTime},
		{"CommitHash", CommitHash},
		{"Description", Description},
	}
	str := ""
	for _, v := range kv {
		str += fmt.Sprintf("%10s: %s\n", v[0], v[1])
	}
	return str
}

// Command 版本命令
func Command() *cobra.Command {
	// versionCmd represents the version command
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf(Print())
		},
	}
	return cmd
}
