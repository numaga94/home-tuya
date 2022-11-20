package lib

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/numaga/home-tuya/utils"
)

func GetDevice(deviceId string) {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, os.Getenv("HOST_URL")+"/v1.0/devices/"+os.Getenv("DEVICE_ID"), bytes.NewReader(body))

	utils.BuildHeader(req, body, Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	log.Println("resp:", string(bs))
}
