# home-tuya

## Build a binary file for ARM devices

- for ARM v6
  ` env GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o main-v6 .

- for ARM v7
  ` env GOOS=linux GOARCH=arm GOARM=7 go build -ldflags="-s -w" -o main-v7 .
