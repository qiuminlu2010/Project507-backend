package main

import (
	"fmt"
	"net/http"

	"qiu/blog/pkg/setting"

	"qiu/blog/model"
	"qiu/blog/pkg/logging"
	"qiu/blog/router"
)

func main() {
	logging.Setup()
	setting.Setup()
	model.Setup()

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
