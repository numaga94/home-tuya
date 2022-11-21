package lib

import (
	"fmt"
	"time"
)

func InOpenHours(beginHour int, endHour int) bool {
	currentHour := time.Now().Hour()
	if currentHour >= beginHour && currentHour <= endHour {
		fmt.Printf("Current time is in open hours between %v and %v.", beginHour, endHour)
		return true
	} else {
		fmt.Printf("Current time is out of the open hours between %v and %v.", beginHour, endHour)
		return false
	}
}
