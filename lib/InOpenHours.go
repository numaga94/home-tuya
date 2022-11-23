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
		fmt.Printf("current time %v:%v is in extended hours.\n", currentHour, currentMinute)
		return true
	} else {
		fmt.Printf("current time %v:%v is out of open hours between %v and %v.\n", currentHour, currentMinute, beginHour, endHour)
		return false
	}
}

func InExtendedHours(beginHour int, endHour int, intervalToUpdateSwitchStatus int) bool {
	currentTime := time.Now()
	currentHour := currentTime.Hour()
	currentMinute := currentTime.Minute()
	if currentHour == endHour+1 && currentMinute <= (59+intervalToUpdateSwitchStatus)%60 {
		fmt.Printf("current time %v:%v is in extended hours.\n", currentHour, currentMinute)
		return true
	} else {
		return false
	}
}
