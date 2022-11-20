package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nuamga/home-tuya/lib"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("loading env failed")
	}
	lib.GetToken()
	lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", true)
	// lib.GetDevice(os.Getenv("DEVICE_ID"))
	time.Sleep(time.Second * 10)
	lib.SwitchDevice(os.Getenv("DEVICE_ID"), "switch_1", false)
	// lib.GetDevice(os.Getenv("DEVICE_ID"))
}
