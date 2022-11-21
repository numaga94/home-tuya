package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func IsOfficeCurrentTempBelowIdealTemp(idealOfficeTemp float64) bool {
	windowTemp := getOfficeTemperature("WINDOW")
	doorTemp := getOfficeTemperature("DOOR")

	averageTemp := (windowTemp + doorTemp) / 2

	if averageTemp > idealOfficeTemp {
		fmt.Println("Current office temperature is at", averageTemp, "degrees, which is above ideal temperature at", idealOfficeTemp, "degrees.")
		return false
	} else {
		fmt.Println("Current office temperature is at", averageTemp, "degrees, which is under ideal temperature at", idealOfficeTemp, "degrees.")
		return true
	}
}

func getOfficeTemperature(location string) float64 {
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
