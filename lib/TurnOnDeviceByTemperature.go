package lib

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/numaga/home-tuya/utils"
)

func TurnOnDeviceByTemperature(idealTemp float64) bool {
	actualTemp := GetCurrentTemperature()

	if int(math.Round(idealTemp)) > int(math.Round(actualTemp)) {
		fmt.Printf("current temperature is at %.1f degrees, which is lower than ideal temperature at %.1f degrees.\n", actualTemp, idealTemp)
		return true
	} else if int(math.Round(idealTemp)) < int(math.Round(actualTemp)) {
		fmt.Printf("current temperature is at %.1f degrees, which is higher than ideal temperature at %.1f degrees.\n", actualTemp, idealTemp)
		return false
	} else {
		fmt.Printf("current temperature is equal to ideal temperature at %.1f degrees.\n", idealTemp)
		return false
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
