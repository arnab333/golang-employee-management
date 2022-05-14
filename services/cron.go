package services

import (
	"context"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func CronInit() {
	c := cron.New()

	// c.AddFunc("0 0 1 */6 *", everySixMonths)
	c.AddFunc("* * * * *", everySixMonths)

	c.Start()

	select {}

	// c.Stop()
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
		_, err := DBConn.findHoliday(context.TODO(), filters)
		if err != nil {
			// fmt.Println("Error findHoliday ==>", err.Error())
			_, err = DBConn.insertHoliday(context.TODO(), val)
			if err != nil {
				fmt.Println("Error insertHoliday ==>", err.Error())
			}
		}
	}
}
