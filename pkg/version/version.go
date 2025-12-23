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
	// CommitMessage is the commit message of the project.
	CommitMessage = "unknown"
)

// PrettyMessage 版本信息
func PrettyMessage() string {
	kv := [][2]string{
		{"Go", runtime.Version()},
		{"Version", Version},
		{"BuildTime", BuildTime},
		{"CommitHash", CommitHash},
		{"CommitMessage", CommitMessage},
	}
	str := ""
	for _, v := range kv {
		str += fmt.Sprintf("%15s: %s\n", v[0], v[1])
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
			cmd.Print(PrettyMessage())
		},
	}
	return cmd
}
