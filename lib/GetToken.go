package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/numaga/home-tuya/utils"
)

var (
	Token        string
	RefreshToken string
	ExpireTime   time.Time
)

type TokenResponse struct {
	Result struct {
		AccessToken  string `json:"access_token"`
		ExpireTime   int    `json:"expire_time"`
		RefreshToken string `json:"refresh_token"`
		UID          string `json:"uid"`
	} `json:"result"`
	Success bool  `json:"success"`
	T       int64 `json:"t"`
}

func GetToken() {
	if Token == "" || RefreshToken == "" || time.Now().After(ExpireTime) {
		InitToken()
	} else {
		TokenRefresh()
	}
}

func InitToken() {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, os.Getenv("HOST_URL")+"/v1.0/token?grant_type=1", bytes.NewReader(body))

	utils.BuildHeader(req, body, Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	ret := TokenResponse{}
	json.Unmarshal(bs, &ret)
	log.Println("resp:", string(bs))

	if ret.Result.AccessToken != "" && ret.Success {
		Token = ret.Result.AccessToken
		RefreshToken = ret.Result.RefreshToken
		ExpireTime = time.Now().Add(time.Second * time.Duration(ret.Result.ExpireTime))
	}
}

func TokenRefresh() {
	method := "GET"
	body := []byte(``)
	req, _ := http.NewRequest(method, fmt.Sprintf("%v/v1.0/token/%v", os.Getenv("HOST_URL"), RefreshToken), bytes.NewReader(body))
	utils.BuildHeader(req, body, RefreshToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	bs, _ := io.ReadAll(resp.Body)
	ret := TokenResponse{}
	json.Unmarshal(bs, &ret)
	log.Println("resp:", string(bs))

	if ret.Result.AccessToken != "" && ret.Success {
		Token = ret.Result.AccessToken
		RefreshToken = ret.Result.RefreshToken
		ExpireTime = time.Now().Add(time.Second * time.Duration(ret.Result.ExpireTime))
	}
}
