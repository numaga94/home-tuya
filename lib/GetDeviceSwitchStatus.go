package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/numaga/home-tuya/utils"
)

type ResponseDeviceStatus struct {
	Result []struct {
		Code  string `json:"code"`
		Value bool   `json:"value"`
	} `json:"result"`
	Success bool   `json:"success"`
	T       int64  `json:"t"`
	Tid     string `json:"tid"`
}

func GetDeviceSwitchStatus(deviceId string) bool {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%v/v1.0/iot-03/devices/%v/status", os.Getenv("HOST_URL"), deviceId), nil)
	utils.BuildHeader(req, nil, Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	deviceStatus := new(ResponseDeviceStatus)
	json.Unmarshal(bs, &deviceStatus)
	for _, v := range deviceStatus.Result {
		if v.Code == os.Getenv("DEVICE_CODE") {
			return v.Value
		}
	}
	fmt.Println(deviceStatus)
	return false
}
