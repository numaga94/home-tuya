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
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
		}

		switch r.Method {
		case "GET":
			responseText := fmt.Sprintf("current ideal office temperature is at %v degrees", idealOfficeTemp)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			temp := r.FormValue("temp")
			idealOfficeTemp, _ = strconv.ParseFloat(temp, 64)
			responseText := fmt.Sprintf("change ideal temperature to %v degrees", idealOfficeTemp)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		default:
			responseText := fmt.Sprintln("sorry, only GET and POST methods are supported.")
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		}
	})

	http.HandleFunc("/ideal-office-hours", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ideal-office-hours" {
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			responseText := fmt.Sprintf("current ideal office hours is between %v and %v.", officeHourBegin, officeHourEnd)
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
		case "POST":
			// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
			begin := r.FormValue("begin")
			end := r.FormValue("end")
			officeHourBegin, _ = strconv.Atoi(begin)
			officeHourEnd, _ = strconv.Atoi(end)
			responseText := fmt.Sprintf("change ideal office hours into between %v and %v", officeHourBegin, officeHourEnd)
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
