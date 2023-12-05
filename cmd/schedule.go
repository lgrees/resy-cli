package cmd

import (
	"github.com/fanniva/resy-cli/internal/schedule"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Schedules a booking",
	Long:  `Schedules a reservation booking to execute at a future time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return schedule.Add(rootCmd.CommandPath())
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
