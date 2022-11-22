# home-tuya

## create a .env file in the root directory for the configs

    # default configs
    HOST_URL="https://openapi.tuyaeu.com"
    CLIENT_ID="YOUR_CLIENT_ID"
    CLIENT_SECRET="YOUR_CLIENT_SECRET"

    # custom configs
    DEVICE_ID="YOUR_DEVICE_ID"
    DEVICE_CODE="YOUR_DEVICE_CODE"
    SENSOR_URLS="http://192.168.1.xxx,http://192.168.1.xxx"
    IDEAL_TEMPERATURE=20
    IDEAL_HUMIDITY=62
    OPEN_HOURS_BEGIN=0
    OPEN_HOURS_END=8
    INTERVAL_TO_CHECK_OPEN_HOURS=1
    INTERVAL_TO_UPDATE_SWITCH_STATUS=5

## build an executable binary file for ARM devices

- for raspberry pi zero: ARM v6

  env GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags="-s -w" -o main-ARMv6 .

- for raspberry pi zero2, 3, and 4: ARM v7

  env GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -ldflags="-s -w" -o main-ARMv7 .
