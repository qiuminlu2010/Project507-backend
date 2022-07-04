package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"

	// "qiu/blog/pkg/logging"
	"qiu/blog/pkg/upload"
)

// @Summary 上传图片
// @Description
// @Tags file
// @Accept multipart/form-data
// @Param image formData file true "image"
// @Param token query string true "token"
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /upload [post]
func UploadImage(c *gin.Context) {

	_, image, err := c.Request.FormFile("image")
	if err != nil {
		// logging.Warn(err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_IMAGE_LOST, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImagePath()
	savePath := upload.GetImagePath()

	src := fullPath + imageName
	// logging.Info("上传图片路径", src, fullPath, savePath)
	// logging.Info("校验图片格式", upload.CheckImageExt(imageName))
	// logging.Info("校验图片大小", upload.CheckImageSize(image))
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(image) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	if err = upload.CheckImage(fullPath); err != nil {
		// logging.Warn(err)
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}
	fmt.Println("保存路径", src)
	if err = c.SaveUploadedFile(image, src); err != nil {
		// logging.Warn(err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}
	thumb_url, err := upload.Thumbnailify(imageName)
	if err != nil {
		fmt.Println(err.Error())
	}
	data := make(map[string]string)
	data["image_url"] = upload.GetImageFullUrl(imageName)
	data["image_save_url"] = savePath + imageName
	data["thumb_url"] = thumb_url
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
