package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func IsCurrentTempUnderIdealTemp(idealOfficeTemp float64) bool {
	urls := getSensorUrlSlice(os.Getenv("SENSOR_URLS"))

	totalTemperature := 0.0
	for _, url := range urls {
		totalTemperature += getDeviceTemperature(url)
	}

	averageTemp := totalTemperature / float64(len(urls))

	if averageTemp > idealOfficeTemp {
		fmt.Println("Current temperature is at", averageTemp, "degrees, which is above ideal temperature at", idealOfficeTemp, "degrees.")
		return false
	} else {
		fmt.Println("Current temperature is at", averageTemp, "degrees, which is under ideal temperature at", idealOfficeTemp, "degrees.")
		return true
	}
}

func getSensorUrlSlice(urls string) []string {
	if strings.Contains(urls, ",") {
		return strings.Split(urls, ",")
	} else {
		return []string{urls}
	}
}

func getDeviceTemperature(url string) float64 {
	requestUrl := fmt.Sprintf("%v/temperature", url)

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
