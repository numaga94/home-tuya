package lib

import (
	"fmt"
	"time"
)

func InOpenHours(beginHour int, endHour int) bool {
	currentHour := time.Now().Hour()
	if currentHour >= beginHour && currentHour <= endHour {
		fmt.Println("Current time is in office hour.")
		return true
	} else {
		fmt.Println("Current time is out of the office hour.")
		return false
	}
}
