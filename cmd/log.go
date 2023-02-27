package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Subcommands related to viewing/removing log files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("log should be run with a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
