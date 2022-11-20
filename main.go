package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/numaga/home-tuya/lib"
)

var (
	Token        string
	RefreshToken string
	ExpireTime   int
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("loading env failed")
	}

	idealOfficeTemp := 30.0
	officeHourBegin := 6
	officeHourEnd := 23
	intervalCheckOfficeHour := 5
	intervalCheckSwitch := 1

	for {
		if !lib.IsInOfficeHour(officeHourBegin, officeHourEnd) {
			time.Sleep(time.Minute * time.Duration(intervalCheckOfficeHour))
		} else {
			// get tuya api token
			lib.GetToken()
			currentDeviceSwitchStatus := lib.GetDeviceSwitchStatus(os.Getenv("DEVICE_ID"))
			isCurrentOfficeTempBelowIdealTemp := lib.IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp)
			// switch office mobile heater by actual office temp
			if isCurrentOfficeTempBelowIdealTemp && !currentDeviceSwitchStatus {
				fmt.Println("Mobile heater is currently off thus turning it on.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", true)
			} else if !isCurrentOfficeTempBelowIdealTemp && currentDeviceSwitchStatus {
				fmt.Println("Mobile heater is currently on thus turning it off.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", false)
			} else {
				fmt.Println("Mobile heater is", currentDeviceSwitchStatus)
			}
			// sleep for 10 minutes
			time.Sleep(time.Minute * time.Duration(intervalCheckSwitch))
		}
	}
}
