package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/setting"
)

func DownloadImg(c *gin.Context) {
	imgType := c.Param("imgType")
	imgName := c.Param("imgName")
	path := ""
	if imgType == "src" {
		path = setting.AppSetting.ImageSavePath
	}
	if imgType == "thumb" {
		path = setting.AppSetting.ThumbSavePath
	}
	if imgType == "temp" {
		path = setting.AppSetting.ImageTempSavePath
	}
	fmt.Println("下载文件", path+"/"+imgName)
	c.FileAttachment(path+"/"+imgName, imgName)
}
