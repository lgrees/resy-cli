package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/book"
	"github.com/bcillie/resy-cli/internal/utils/paths"
	"github.com/rs/zerolog"
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
		flags := cmd.Flags()

		venueId, _ := flags.GetInt32("venueId")
		partySize, _ := flags.GetString("partySize")
		reservationDate, _ := flags.GetString("reservationDate")
		bookingDateTime, _ := flags.GetString("bookingDateTime")
		reservationTimes, _ := flags.GetStringSlice("reservationTimes")
		reservationTypes, _ := flags.GetStringSlice("reservationTypes")
		dryRun, _ := flags.GetBool("dryRun")

		bookingDetails := &book.BookingDetails{
			VenueId:          string(venueId),
			PartySize:        partySize,
			BookingDateTime:  bookingDateTime,
			ReservationDate:  reservationDate,
			ReservationTimes: reservationTimes,
			ReservationTypes: reservationTypes,
		}

		p, err := paths.GetAppPaths()
		if err != nil {
			return err
		}

		venueDetails, _ := api.GetConfig(venueId)
		formattedTime := time.Now().Format("Mon Jan _2 15:04:05 2006")

		logFileName := path.Join(p.LogPath, fmt.Sprintf("%s_%s.log", venueDetails.Name, formattedTime))
		logFile, err := os.OpenFile(
			logFileName,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}

		defer logFile.Close()

		l := zerolog.New(logFile).With().Timestamp().Logger()

		l.Info().Object("booking_details", bookingDetails).Msg("starting book job")

		if bookingDateTime != "" {
			err = book.WaitThenBook(bookingDetails, dryRun, l)
		} else {
			err = book.Book(bookingDetails, dryRun, l)
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(bookCmd)

	flags := bookCmd.Flags()

	flags.String("venueId", "", "The venue id of the restaurant")
	flags.Bool("dryRun", false, "When true, skips booking")
	flags.String("partySize", "", "The party size for the reservation")
	flags.String("bookingDateTime", "", "The time when the reservation should be booked")
	flags.String("reservationDate", "", "The date of the reservation")
	flags.StringSlice("reservationTimes", make([]string, 0), "The times for the reservation")
	flags.StringSlice("reservationTypes", make([]string, 0), "The table types for the reservation")

	bookCmd.MarkFlagRequired("venueId")
	bookCmd.MarkFlagRequired("partySize")
	bookCmd.MarkFlagRequired("reservationDate")
	bookCmd.MarkFlagRequired("reservationTimes")
}
