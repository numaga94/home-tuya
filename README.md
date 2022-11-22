# home-tuya

## Env

`# default configs
HOST_URL="https://openapi.tuyaeu.com"
CLIENT_ID=
CLIENT_SECRET=

# custom configs

DEVICE_ID=
DEVICE_CODE="switch"
SENSOR_URLS="http://192.168.1.xxx,http://192.168.1.xxx"
IDEAL_TEMPERATURE=20
IDEAL_HUMIDITY=62
OPEN_HOURS_BEGIN=0
OPEN_HOURS_END=8
INTERVAL_TO_CHECK_OPEN_HOURS=1
INTERVAL_TO_UPDATE_SWITCH_STATUS=5`

## Build a binary file for ARM devices

- for ARM v6
  `env GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o main-v6 .`

- for ARM v7
  `env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o main-v7 .`
