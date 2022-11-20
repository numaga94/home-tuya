package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/nuamga/home-tuya/lib"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("loading env failed")
	}

	idealOfficeTemp := 28.0

	for {
		// get tuya api token
		lib.GetToken()
		// switch office mobile heater by actual office temp
		if IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp) {
			lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", true)
		} else {
			lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", false)
		}
		// sleep for 10 minutes
		time.Sleep(time.Minute * 2)
	}
}

func IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp float64) bool {
	windowTemp := GetOfficeTemperature("WINDOW")
	doorTemp := GetOfficeTemperature("DOOR")

	averageTemp := (windowTemp + doorTemp) / 2

	if averageTemp > idealOfficeTemp {
		fmt.Println("current office temperature is at", averageTemp, "degrees, which is above ideal temperature of", idealOfficeTemp, "degrees. Turn off office mobile heater.")
		return false
	} else {
		fmt.Println("current office temperature is at", averageTemp, "degrees, which is blow ideal temperature of", idealOfficeTemp, "degrees. Turn on office mobile heater.")
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
