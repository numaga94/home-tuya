package main

import (
	"fmt"
	"log"
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

	idealOfficeTemp := 28.0
	officeHourBegin := 6
	officeHourEnd := 23
	intervalCheckOfficeHour := 1
	intervalCheckSwitch := 10

	go func() {
		for {
			if !lib.IsInOfficeHour(officeHourBegin, officeHourEnd) {
				time.Sleep(time.Minute * time.Duration(intervalCheckOfficeHour))
			} else {
				// get tuya api token
				if err := lib.GetToken(); err != nil {
					fmt.Println(err.Error())
					time.Sleep(time.Minute * time.Duration(intervalCheckOfficeHour))
					continue
				} else {
					// get current state
					currentDeviceSwitchStatus := lib.GetDeviceSwitchStatus(os.Getenv("DEVICE_ID"))
					isCurrentOfficeTempBelowIdealTemp := lib.IsOfficeCurrentTempUnderIdealTemp(idealOfficeTemp)
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
	}()

	http.HandleFunc("/ideal-office-temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ideal-office-temperature" {
			http.Error(w, "not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			fmt.Fprintf(w, "current ideal office temperature is at %v degrees", idealOfficeTemp)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			temp := r.FormValue("temp")
			idealOfficeTemp, _ = strconv.ParseFloat(temp, 64)
			fmt.Fprintf(w, "change ideal temperature to %v degrees", temp)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	})
	fmt.Println("http server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
