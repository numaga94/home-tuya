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

	idealTemperature, _ := strconv.ParseFloat(os.Getenv("IDEAL_TEMPERATURE"), 64)
	idealHumidity, _ := strconv.ParseFloat(os.Getenv("IDEAL_HUMIDITY"), 64)
	openHoursBegin, _ := strconv.Atoi(os.Getenv("OPEN_HOURS_BEGIN"))
	openHoursEnd, _ := strconv.Atoi(os.Getenv("OPEN_HOURS_END"))
	intervalToCheckOpenHours, _ := strconv.Atoi(os.Getenv("INTERVAL_TO_CHECK_OPEN_HOURS"))
	intervalToUpdateSwitchStatus, _ := strconv.Atoi(os.Getenv("INTERVAL_TO_UPDATE_SWITCH_STATUS"))

	go func() {
		for {
			// check if current hour in open hours
			if !lib.InOpenHours(openHoursBegin, openHoursEnd) {
				time.Sleep(time.Minute * time.Duration(intervalToCheckOpenHours))
				continue
			}

			// get tuya api token
			if err := lib.GetToken(); err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Minute * time.Duration(intervalToCheckOpenHours))
				continue
			}

			// get current state
			currentDeviceSwitchStatus := lib.GetDeviceSwitchStatus(os.Getenv("DEVICE_ID"))
			isCurrentTempUnderIdealTemp := lib.IsCurrentTempUnderIdealTemp(idealTemperature)
			// fmt.Println("current device switch status:", currentDeviceSwitchStatus)
			// fmt.Println("is current temp under ideal temperature:", isCurrentTempUnderIdealTemp)
			// switch office mobile heater by actual office temp
			if isCurrentTempUnderIdealTemp && !currentDeviceSwitchStatus {
				fmt.Println("Mobile heater is currently off thus turning it on.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), true)
			} else if !isCurrentTempUnderIdealTemp && currentDeviceSwitchStatus {
				fmt.Println("Mobile heater is currently on thus turning it off.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), false)
			} else {
				fmt.Println("Mobile heater is", currentDeviceSwitchStatus)
			}
			// sleep for 10 minutes
			time.Sleep(time.Minute * time.Duration(intervalToUpdateSwitchStatus))
		}
	}()

	http.HandleFunc("/ideal-temperature-humidity", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ideal-temperature-humidity" {
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
		}

		switch r.Method {
		case "GET":
			currentTemperature := lib.GetCurrentTemperature()
			currentHumidity := lib.GetCurrentHumidity()
			responseText := fmt.Sprintf("ideal: %v degrees + %v %%H and current: %v degrees + %v %%H", idealTemperature, idealHumidity, currentTemperature, currentHumidity)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// parse values and reassign them to global variants
			temp := r.FormValue("temperature")
			humidity := r.FormValue("humidity")
			if idealT, err := strconv.ParseFloat(temp, 64); err == nil || idealT == 0.0 {
				idealTemperature = idealT
			}
			if idealH, err := strconv.ParseFloat(humidity, 64); err == nil || idealH == 0.0 {
				idealHumidity = idealH
			}
			responseText := fmt.Sprintf("change ideal temperature to %v degrees and ideal humidity to %v %%H", idealTemperature, idealHumidity)
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
