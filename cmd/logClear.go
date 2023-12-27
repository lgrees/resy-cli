package cmd

import (
	"github.com/bcillie/resy-cli/internal/log"
	"github.com/spf13/cobra"
)

var logClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all log files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return log.Clear()
	},
}

func init() {
	logCmd.AddCommand(logClearCmd)
}
