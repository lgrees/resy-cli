package cmd

import (
	"github.com/lgrees/resy-cli/internal/book"
	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "(internal) Books a reservation immediately",
	Long: `
	Books a reservation using the resy API. This command exists for internal use.
	Generally, users of resy-cli should schedule a booking using "resy schedule".
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// var bookingDetails book.BookingDetails

		flags := cmd.Flags()

		venueId, _ := flags.GetString("venueId")
		partySize, _ := flags.GetString("partySize")
		reservationDate, _ := flags.GetString("reservationDate")
		bookingDateTime, _ := flags.GetString("bookingDateTime")
		reservationTimes, _ := flags.GetStringSlice("reservationTimes")
		reservationTypes, _ := flags.GetStringSlice("reservationTypes")
		dryRun, _ := flags.GetBool("dryRun")
		wait, _ := flags.GetBool("wait")

		bookingDetails := &book.BookingDetails{
			VenueId:          venueId,
			PartySize:        partySize,
			BookingDateTime:  bookingDateTime,
			ReservationDate:  reservationDate,
			ReservationTimes: reservationTimes,
			ReservationTypes: reservationTypes,
		}

		var err error

		if wait {
			err = book.WaitThenBook(bookingDetails, dryRun)
		} else {
			err = book.Book(bookingDetails, dryRun)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)

	flags := bookCmd.Flags()

	flags.String("venueId", "", "The venue id of the restaurant")
	flags.Bool("dryRun", false, "When true, skips booking")
	flags.Bool("wait", true, "When true, waits for bookingDateTime to book")
	flags.String("partySize", "", "The party size for the reservation")
	flags.String("bookingDateTime", "", "The time when the reservation should be booked")
	flags.String("reservationDate", "", "The date of the reservation")
	flags.StringSlice("reservationTimes", make([]string, 0), "The times for the reservation")
	flags.StringSlice("reservationTypes", make([]string, 0), "The table types for the reservation")

	bookCmd.MarkFlagRequired("venueId")
	bookCmd.MarkFlagRequired("partySize")
	bookCmd.MarkFlagRequired("reservationDate")
	bookCmd.MarkFlagRequired("reservationTimes")
	bookCmd.MarkFlagRequired("reservationTypes")
	bookCmd.MarkFlagsRequiredTogether("wait", "bookingDateTime")
}
