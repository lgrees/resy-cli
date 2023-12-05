package cmd

import (
	"github.com/fanniva/resy-cli/internal/ping"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping the resy API",
	Long:  `Ping the resy API to verify that the correct credentials are stored.`,
	Run: func(cmd *cobra.Command, args []string) {
		ping.Ping()
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
