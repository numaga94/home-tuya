package utils

import "strings"

func GetSensorUrlSlice(urls string) []string {
	if strings.Contains(urls, ",") {
		return strings.Split(urls, ",")
	} else {
		return []string{urls}
	}
}
