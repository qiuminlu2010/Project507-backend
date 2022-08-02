package v1

import (
	"github.com/gin-gonic/gin"
)

func DownloadImg(c *gin.Context) {
	// imgType := c.Param("imgType")
	// imgName := c.Param("imgName")
	// path := ""
	// if imgType == "src" {
	// 	path = setting.MinioSetting.ImageSavePath
	// }
	// if imgType == "thumb" {
	// 	path = setting.MinioSetting.ThumbSavePath
	// }
	// if imgType == "temp" {
	// 	path = setting.MinioSetting.ImageTempSavePath
	// }
	// if imgType == "avatar" {
	// 	path = setting.MinioSetting.AvatarSavePath
	// }
	// fmt.Println("下载文件", path+imgName, setting.MinioSetting.Host+path+imgName)
	// c.FileAttachment(setting.MinioSetting.Host+path+imgName, imgName)
}
