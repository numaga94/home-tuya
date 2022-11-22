package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/numaga/home-tuya/utils"
)

func IsCurrentHumidityUnderIdealHumidity(idealHumidity float64) bool {
	averageHumidity := GetCurrentHumidity()
	if averageHumidity < idealHumidity {
		fmt.Println("Current humidity is at", averageHumidity, "degrees, which is under ideal humidity at", idealHumidity, "%H.")
		return true
	} else {
		fmt.Println("Current humidity is at", averageHumidity, "degrees, which is above ideal humidity at", idealHumidity, "%H.")
		return false
	}
}

func GetCurrentHumidity() float64 {
	urls := utils.GetSensorUrlSlice(os.Getenv("SENSOR_URLS"))

	totalHumidity := 0.0
	for _, url := range urls {
		totalHumidity += getDeviceHumidity(url)
	}

	return totalHumidity / float64(len(urls))
}

func getDeviceHumidity(url string) float64 {
	requestUrl := fmt.Sprintf("%v/humidity", url)

	req, _ := http.NewRequest("GET", requestUrl, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	humidity, err := strconv.ParseFloat(string(bs), 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	return humidity
}
