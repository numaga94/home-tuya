package lib

import (
	"fmt"
	"time"
)

func InOpenHours(beginHour int, endHour int) bool {
	currentHour := time.Now().Hour()
	if currentHour >= beginHour && currentHour <= endHour {
		fmt.Println("Current time is in open hours.")
		return true
	} else {
		fmt.Println("Current time is out of the open hours.")
		return false
	}
}
