package cmd

import (
	"github.com/lgrees/resy-cli/internal/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize resy credentials (required for new users)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setup.SurveyConfig()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
