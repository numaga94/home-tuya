package lib

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/nuamga/home-tuya/utils"
)

func SwitchDevice(deviceId string, code string, value bool) {
	method := "POST"
	body := []byte(fmt.Sprintf(`{
					"commands": [
						{
						"code": "%v",
						"value": %v
						}
					]
					}`, code, value))
	req, _ := http.NewRequest(method, os.Getenv("HOST_URL")+"/v1.0/iot-03/devices/"+os.Getenv("DEVICE_ID")+"/commands", bytes.NewReader(body))

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
