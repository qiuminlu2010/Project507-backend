package cron

import (
	log "qiu/blog/pkg/logging"
	"time"
)

func ClearLoggingFile() {
	time.Sleep(time.Second)
	log.Info("clearing")
}
