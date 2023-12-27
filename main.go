package main

import (
	"fmt"
	"time"

	"github.com/lgrees/resy-cli/cmd"
	"github.com/lgrees/resy-cli/internal/api"
	"github.com/lgrees/resy-cli/internal/utils"
)

type TestS struct {
	Name string    `query:"name"`
	Id   int32     `query:"venue_id"`
	Date time.Time `query:"time" fmt:"Mon, 02 Jan 2006 15:04:05 MST"`
}

func main() {
	cmd.Execute()
	sv, err := api.SearchVenues("Misi")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("%s", sv)

	s := TestS{"Ben", 3, time.Now()}
	utils.GetQueryParams(s)

	t, _ := time.Parse(time.RFC822, "06 Jan 24 18:00 EDT")
	f, err := api.Find(&api.FindParams{VenueId: 1010, PartySize: 2, ReservationDate: t})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("%s", f)

}
