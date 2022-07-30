package v1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"qiu/backend/pkg/e"
	gin_http "qiu/backend/pkg/http"
	log "qiu/backend/pkg/logging"

	// "qiu/backend/pkg/logging"

	"qiu/backend/pkg/minio"
	"qiu/backend/pkg/setting"
	"qiu/backend/pkg/upload"
)

// @Summary 上传图片
// @Description
// @Tags file
// @Accept multipart/form-data
// @Param image formData file true "image"
// @Param token header string true "token"
// @Produce  json
// @Router /api/v1/upload/image [post]
func UploadImage(c *gin.Context) {

	_, image, err := c.Request.FormFile("image")
	if err != nil {
		log.Logger.Error("保存图片失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_IMAGE_LOST, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	// fullPath := upload.GetImagePath()
	src := "/" + setting.MinioSetting.TempBucketName + "/" + imageName

	// src := savePath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(image) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}
	if err = minio.PutImage(setting.MinioSetting.TempBucketName, imageName, image); err != nil {
		log.Logger.Error("保存图片失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	// if err = c.SaveUploadedFile(image, setting.MinioSetting.Host+src); err != nil {
	// 	log.Logger.Error("保存图片失败", err)
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
	// 	return
	// }
	log.Logger.Info("保存上传图片", src)
	data := make(map[string]string)
	data["image_url"] = src

	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 上传视频
// @Description
// @Tags file
// @Accept multipart/form-data
// @Param video formData file true "video"
// @Param token header string true "token"
// @Produce  json
// @Router /api/v1/upload/video [post]
func UploadVideo(c *gin.Context) {

	_, video, err := c.Request.FormFile("video")
	if err != nil {
		log.Logger.Error("保存视频失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if video == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_IMAGE_LOST, nil)
		return
	}
	fileName := upload.GetFileName(video.Filename)
	// videoName := upload.GetImageName(video.Filename)
	// fullPath := upload.GetImagePath()
	// videoSavePath := upload.GetVideoPath()
	// previewSavePath := upload.GetVideoPreviewPath()

	videoSrc := "/" + setting.MinioSetting.VideoBucketName + "/" + fileName + ".ts"
	preiviewSrc := "/" + setting.MinioSetting.PreviewBucketName + "/" + fileName + ".jpg"

	if !upload.CheckVideoSize(video) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}
	// TODO: 保存到minio,通过消息队列异步压缩视频和生成预览图
	tempSavePath := "runtime/temp/" + fileName
	if err = c.SaveUploadedFile(video, tempSavePath); err != nil {
		log.Logger.Error("保存文件失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}
	// minio.PutVideoAndPreview(video, fileName)
	if err = minio.PutVideoAndPreview(tempSavePath, fileName); err != nil {
		log.Logger.Error("保存视频失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		os.Remove(tempSavePath)
		return
	}
	os.Remove(tempSavePath)
	log.Logger.Info("保存上传视频", videoSrc)
	data := make(map[string]string)
	data["video_url"] = videoSrc
	data["preview_url"] = preiviewSrc
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
