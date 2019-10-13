package main

import (
	"strconv"
	"time"
)

func currentTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
