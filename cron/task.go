package cron

import (
	"fmt"
	"time"
)

func ClearLoggingFile() {
	time.Sleep(time.Second)
	fmt.Println("clearing")
}

