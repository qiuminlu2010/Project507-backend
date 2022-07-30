package cron

import (
	log "qiu/backend/pkg/logging"
	"time"
)

func ClearLoggingFile() {
	time.Sleep(time.Second)
	log.Logger.Info("clearing")
}
