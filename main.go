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

	idealTemperature := 28.0
	openHoursBegin := 6
	openHoursEnd := 23
	intervalToCheckOpenHours := 1
	intervalToUpdateSwitchStatus := 10

	go func() {
		for {
			if !lib.InOpenHours(openHoursBegin, openHoursEnd) {
				time.Sleep(time.Minute * time.Duration(intervalToCheckOpenHours))
			} else {
				// get tuya api token
				if err := lib.GetToken(); err != nil {
					fmt.Println(err.Error())
					time.Sleep(time.Minute * time.Duration(intervalToCheckOpenHours))
					continue
				} else {
					// get current state
					currentDeviceSwitchStatus := lib.GetDeviceSwitchStatus(os.Getenv("DEVICE_ID"))
					isCurrentTempBelowIdealTemp := lib.IsCurrentTempUnderIdealTemp(idealTemperature)
					// switch office mobile heater by actual office temp
					if isCurrentTempBelowIdealTemp && !currentDeviceSwitchStatus {
						fmt.Println("Mobile heater is currently off thus turning it on.")
						lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", true)
					} else if !isCurrentTempBelowIdealTemp && currentDeviceSwitchStatus {
						fmt.Println("Mobile heater is currently on thus turning it off.")
						lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", false)
					} else {
						fmt.Println("Mobile heater is", currentDeviceSwitchStatus)
					}
					// sleep for 10 minutes
					time.Sleep(time.Minute * time.Duration(intervalToUpdateSwitchStatus))
				}
			}
		}
	}()

	http.HandleFunc("/ideal-temperature", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ideal-temperature" {
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
		}

		switch r.Method {
		case "GET":
			responseText := fmt.Sprintf("current ideal temperature is at %v degrees", idealTemperature)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// parse values and reassign them to global variants
			temp := r.FormValue("temp")
			idealTemperature, _ = strconv.ParseFloat(temp, 64)
			responseText := fmt.Sprintf("change ideal temperature to %v degrees", idealTemperature)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		default:
			responseText := fmt.Sprintln("sorry, only GET and POST methods are supported.")
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		}
	})

	http.HandleFunc("/open-hours", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-hours" {
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			responseText := fmt.Sprintf("current open hours is between %v and %v.", openHoursBegin, openHoursEnd)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// parse values and reassign them to global variants
			begin := r.FormValue("begin")
			end := r.FormValue("end")
			openHoursBegin, _ = strconv.Atoi(begin)
			openHoursEnd, _ = strconv.Atoi(end)
			responseText := fmt.Sprintf("change opens hours into between %v and %v.", openHoursBegin, openHoursEnd)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		default:
			responseText := fmt.Sprintln("sorry, only GET and POST methods are supported.")
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		}
	})

	fmt.Println("http server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
