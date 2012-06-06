package ibtws

import (
	"strconv"
	"time"
)

const ibTime = "20060102 15:04:05 MST"

func decodeServerTime(data string) (time.Time, error) {
	return time.Parse(ibTime, data)
}

func decodeServerVersion(data string) (int, error) {
	return strconv.Atoi(data)
}

func encodeClientVersion(ver int) string {
	return strconv.Itoa(ver)
}
