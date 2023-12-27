package cmd

import (
	"github.com/bcillie/resy-cli/internal/log"
	"github.com/spf13/cobra"
)

var logViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View log files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return log.View()
	},
}

func init() {
	logCmd.AddCommand(logViewCmd)
}
