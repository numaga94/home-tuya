package lib

import (
	"fmt"
	"time"
)

func InOpenHours(beginHour int, endHour int, intervalToUpdateSwitchStatus int) bool {
	currentTime := time.Now()
	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	if currentHour >= beginHour && currentHour <= endHour {
		fmt.Println("current time is in open hours between", beginHour, "and", endHour)
		return true
	} else if currentHour == endHour+1 && currentMinute <= (59+intervalToUpdateSwitchStatus)%60 {
		fmt.Println("current time is in open hours between", beginHour, "and", endHour)
		return true
	} else {
		fmt.Println("current time is out of the open hours between", beginHour, "and", endHour)
		return false
	}
}
