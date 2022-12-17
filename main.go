package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	// set default SWITCH to ON
	SWITCH := true

	go func() {
		for {
			// check if current hour in open hours
			if !lib.InOpenHours(openHoursBegin, openHoursEnd, intervalToUpdateSwitchStatus) {
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
			turnOnDeviceByTemperature := lib.TurnOnDeviceByTemperature(idealTemperature)
			turnOnDeviceByHumidity := lib.TurnOnDeviceByHumidity(idealHumidity)
			fmt.Println("Turn on device by temperature:", turnOnDeviceByTemperature, "| Turn on device by humidity:", turnOnDeviceByHumidity, "| Current device status:", currentDeviceSwitchStatus)

			// check if switch is ON
			if !SWITCH {
				if currentDeviceSwitchStatus {
					fmt.Println("SWITCH is OFF thus turning it off")
					lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), false)
				} else {
					fmt.Println("SWITCH is OFF and heater is OFF")
				}
				// switch office mobile heater by actual office temp
			} else if lib.InExtendedHours(openHoursBegin, openHoursEnd, intervalToUpdateSwitchStatus) {
				// action on switch
				fmt.Println("mobile heater is in extended hours thus turning it off.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), false)
			} else if turnOnDeviceByTemperature || turnOnDeviceByHumidity {
				if !currentDeviceSwitchStatus {
					// action on switch
					fmt.Println("mobile heater is currently off thus turning it on.")
					lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), true)
				} else {
					fmt.Println("mobile heater is currently on.")
				}
			} else if currentDeviceSwitchStatus {
				// action on switch
				fmt.Println("mobile heater is currently on thus turning it off.")
				lib.SwitchDevice(os.Getenv("DEVICE_ID"), os.Getenv("DEVICE_CODE"), false)
			} else {
				fmt.Println("mobile heater is currently off.")
			}

			// sleep for 10 minutes for next loop
			time.Sleep(time.Minute * time.Duration(intervalToUpdateSwitchStatus))
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// if r.URL.Path != "/ideal-temperature-humidity" {
		// 	responseText := fmt.Sprintln("404, not found.")
		// 	log.Println(responseText)
		// 	http.Error(w, responseText, http.StatusNotFound)
		// }
		switchStatus := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("switch")))
		idealTemperatureHumidity := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("ideal-temperature-humidity")))
		openHours := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("open-hours")))
		currentTemperatureHumidity := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("current-temperature-humidity")))

		if switchStatus == "true" {
			SWITCH = true
			responseText := "heater SWITCH turns ON"
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
			return
		} else if switchStatus == "false" {
			SWITCH = false
			responseText := "heater SWITCH turns OFF"
			log.Println(responseText)
			fmt.Fprintln(w, responseText)
			return
		} else if currentTemperatureHumidity == "true" {
			currentTemperature := lib.GetCurrentTemperature()
			currentHumidity := lib.GetCurrentHumidity()

			log.Println("current temperature:", currentTemperature, "current humidity:", currentHumidity)
			response, _ := json.Marshal(map[string]float64{
				"temperature": currentTemperature,
				"humidity":    currentHumidity,
			})
			fmt.Fprintln(w, string(response))
			return
		} else if idealTemperatureHumidity == "true" {
			switch r.Method {
			case "GET":
				currentTemperature := lib.GetCurrentTemperature()
				currentHumidity := lib.GetCurrentHumidity()
				responseText := fmt.Sprintf("ideal: %.1f degrees, %.1f %%H and current: %.1f degrees, %.1f %%H", idealTemperature, idealHumidity, currentTemperature, currentHumidity)
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
				if idealT, err := strconv.ParseFloat(temp, 64); err == nil {
					idealTemperature = idealT
				}
				if idealH, err := strconv.ParseFloat(humidity, 64); err == nil {
					idealHumidity = idealH
				}
				responseText := fmt.Sprintf("change ideal temperature to %.1f degrees and ideal humidity to %.1f %%H", idealTemperature, idealHumidity)
				log.Println(responseText)
				fmt.Fprintln(w, responseText)
			default:
				responseText := fmt.Sprintln("sorry, only GET and POST methods are supported.")
				log.Println(responseText)
				fmt.Fprintln(w, responseText)
			}
			return
		} else if openHours == "true" {
			switch r.Method {
			case "GET":
				responseText := fmt.Sprintf("current open hours is between %d and %d.", openHoursBegin, openHoursEnd)
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
				responseText := fmt.Sprintf("change opens hours into between %d and %d.", openHoursBegin, openHoursEnd)
				log.Println(responseText)
				fmt.Fprintln(w, responseText)
			default:
				responseText := fmt.Sprintln("sorry, only GET and POST methods are supported.")
				log.Println(responseText)
				fmt.Fprintln(w, responseText)
			}
			return
		} else {
			responseText := fmt.Sprintln("404, not found.")
			log.Println(responseText)
			http.Error(w, responseText, http.StatusNotFound)
			return
		}
	})

	fmt.Println("http server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
