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

func TurnOnDeviceByHumidity(idealHumidity float64) bool {
	actualHumidity := GetCurrentHumidity()

	if int(math.Round(idealHumidity)) > int(math.Round(actualHumidity)) {
		fmt.Printf("current humidity is at %.1f %%H, which feels drier than ideal humidity at %.1f %%H.\n", actualHumidity, idealHumidity)
		return false
	} else if int(math.Round(idealHumidity)) < int(math.Round(actualHumidity)) {
		fmt.Printf("current humidity is at %.1f %%H, which feels wetter than ideal humidity at %.1f %%H.\n", actualHumidity, idealHumidity)
		return true
	} else {
		fmt.Printf("current humidity is equal to ideal humidity at %.1f %%H.\n", actualHumidity)
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
