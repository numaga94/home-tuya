package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/numaga/home-tuya/utils"
)

func IsCurrentTempUnderIdealTemp(idealTemp float64) bool {
	averageTemp := GetCurrentTemperature()

	if averageTemp > idealTemp {
		fmt.Println("Current temperature is at", averageTemp, "degrees, which is above ideal temperature at", idealTemp, "degrees.")
		return false
	} else {
		fmt.Println("Current temperature is at", averageTemp, "degrees, which is under ideal temperature at", idealTemp, "degrees.")
		return true
	}
}

func GetCurrentTemperature() float64 {
	urls := utils.GetSensorUrlSlice(os.Getenv("SENSOR_URLS"))

	totalTemperature := 0.0
	for _, url := range urls {
		totalTemperature += getDeviceTemperature(url)
	}

	return totalTemperature / float64(len(urls))
}

func getDeviceTemperature(url string) float64 {
	requestUrl := fmt.Sprintf("%v/temperature", url)

	req, _ := http.NewRequest("GET", requestUrl, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	temp, err := strconv.ParseFloat(string(bs), 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	return temp
}
