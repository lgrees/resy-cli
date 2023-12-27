package main

import (
	"fmt"
	"time"

	"github.com/bcillie/resy-cli/cmd"
	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/utils/date"
)

func main() {
	cmd.Execute()
	// sv, err := api.SearchVenues("Misi")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println("%s", sv)

	t, _ := (time.Parse(time.RFC822, "06 Jan 24 18:00 EDT"))
	rd := *date.NewResyDate(t, time.DateOnly)
	f, err := api.Find(&api.FindParams{VenueId: 1010, PartySize: 2, ReservationDate: rd})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dp := api.DetailsParams{ConfigId: f[0].Config.Token, Day: rd, PartySize: 2}
	fmt.Printf("%s", dp)

}
