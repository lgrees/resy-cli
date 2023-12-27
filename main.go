package main

import (
	"github.com/bcillie/resy-cli/cmd"
)

func main() {
	cmd.Execute()
	// sv, err := api.SearchVenues("Misi")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println("%s", sv)

	// t, _ := (time.Parse(time.RFC822, "12 Jan 24 18:00 EDT"))
	// rd, _ := date.NewResyDate(t, time.DateOnly)
	// _, err := api.Find(&api.FindParams{VenueId: 1010, PartySize: 3, ReservationDate: *rd})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// dp := api.DetailsParams{ConfigId: f[0].Config.Token, Day: *rd, PartySize: 2}
	// _, err = api.GetDetails(&dp)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// pp, err := json.MarshalIndent(dr, "", "  ")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Printf("%s", string(pp))

	// err = api.Book(dr)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

}
