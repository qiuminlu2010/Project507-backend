package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/setting"

	"qiu/blog/cron"
	"qiu/blog/model"
	"qiu/blog/pkg/logging"
	"qiu/blog/router"
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
		fmt.Println("结束程序")
		os.Exit(1)
	}()

	logging.Setup()
	setting.Setup()
	model.Setup()
	redis.Setup()
	cron.Setup()

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

	s.ListenAndServe()
}
