package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/numaga/home-tuya/lib"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("loading env failed")
	}

	idealOfficeTemp := 30.0
	officeHourBegin := 6
	officeHourEnd := 23
	intervalCheckOfficeHour := 5
	intervalCheckSwitch := 10

	for {
		if !IsInOfficeHour(officeHourBegin, officeHourEnd) {
			time.Sleep(time.Minute * time.Duration(intervalCheckOfficeHour))
		} else {
			// get tuya api token
			lib.GetToken()
			currentDeviceSwitchStatus := lib.GetDeviceSwitchStatus(os.Getenv("DEVICE_ID"))
			isCurrentOfficeTempBelowIdealTemp := IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp)
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

func IsInOfficeHour(beginHour int, endHour int) bool {
	currentHour := time.Now().Hour()
	if currentHour >= beginHour && currentHour <= endHour {
		fmt.Println("Current time is in office hour.")
		return true
	} else {
		fmt.Println("Current time is out of the office hour.")
		return false
	}
}

func IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp float64) bool {
	windowTemp := GetOfficeTemperature("WINDOW")
	doorTemp := GetOfficeTemperature("DOOR")

	averageTemp := (windowTemp + doorTemp) / 2

	if averageTemp > idealOfficeTemp {
		fmt.Println("Current office temperature is at", averageTemp, "degrees, which is above ideal temperature of", idealOfficeTemp, "degrees.")
		return false
	} else {
		fmt.Println("Current office temperature is at", averageTemp, "degrees, which is blow ideal temperature of", idealOfficeTemp, "degrees.")
		return true
	}
}

func GetOfficeTemperature(location string) float64 {
	var requestUrl string
	if location == "WINDOW" {
		requestUrl = fmt.Sprintf("%v/temperature", os.Getenv("SENSOR_WINDOW_URL"))
	} else {
		requestUrl = fmt.Sprintf("%v/temperature", os.Getenv("SENSOR_DOOR_URL"))
	}

	req, _ := http.NewRequest("GET", requestUrl, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	temp, err := strconv.ParseFloat(string(bs), 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	return temp
}
