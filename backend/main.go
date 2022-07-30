package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"qiu/backend/cron"
	"qiu/backend/model"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/minio"
	"qiu/backend/pkg/redis"
	"qiu/backend/pkg/setting"
	"qiu/backend/router"
	msg "qiu/backend/service/msg"
	"syscall"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		cron.Exit()
		log.Logger.Info("服务关闭")
		os.Exit(1)
	}()

	log.Setup()
	setting.Setup()
	model.Setup()
	redis.Setup()
	minio.Setup()
	cron.Setup()
	msg.Setup()
	// service.FlushArticleLikeUsers()
	/*
		router := gin.Default()
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "test",
			})
		})
	*/
	router := router.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	log.Logger.Info("服务启动")
	s.ListenAndServe()
}
