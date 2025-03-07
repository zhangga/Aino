package main

import (
	"github.com/spf13/cobra"
	"github.com/zhangga/aino/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:     "aino",
	Short:   "aino service",
	Version: version.Version,
}

func init() {
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(cmdRun)
}
