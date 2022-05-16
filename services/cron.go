package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func CronInit() {
	c := cron.New()

	c.AddFunc("0 0 1 */6 *", everySixMonths) // This will run every six months
	// c.AddFunc("* * * * *", everySixMonths) // This will run every one minute

	c.Start()

	for {
		select {}
	}

	// go func() {
	// 	for {
	// 		everySixMonths()
	// 		<-time.After(1 * time.Minute)
	// 	}
	// }()

}

func everySixMonths() {
	fmt.Println("Every 6 months")

	holidays := calendarInit()

	if holidays == nil {
		fmt.Println("Sorry! No Holidays Found.")
		return
	}

	for _, val := range holidays {
		filters := bson.M{
			"date": val.Date,
		}
		_, err := DBConn.FindHoliday(context.TODO(), filters)
		if err != nil {
			// fmt.Println("Error findHoliday ==>", err.Error())
			val.Summary = strings.ReplaceAll(val.Summary, "(Regional Holiday)", "")
			val.Summary = strings.ReplaceAll(val.Summary, "West Bengal:", "")
			val.Summary = strings.TrimSpace(val.Summary)
			val.IsActive = true
			_, err = DBConn.InsertHoliday(context.TODO(), val)
			if err != nil {
				fmt.Println("Error insertHoliday ==>", err.Error())
			}
		}
	}
}
