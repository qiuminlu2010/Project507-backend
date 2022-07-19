package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	log "qiu/blog/pkg/logging"

	// "qiu/blog/pkg/logging"
	"qiu/blog/pkg/upload"
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
	savePath := upload.GetImageTempPath()

	src := savePath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(image) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	if err = c.SaveUploadedFile(image, src); err != nil {
		log.Logger.Error("保存图片失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}
	log.Logger.Error("保存上传图片", src)
	data := make(map[string]string)
	data["image_url"] = src

	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
