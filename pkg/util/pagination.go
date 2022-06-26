package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"qiu/blog/pkg/setting"
)

func GetPage(c *gin.Context) (int, int) {
	result := 0
	page := 0
	page, _ = com.StrTo(c.Query("pageNum")).Int()
	if page > 0 {
		result = (page) * setting.AppSetting.PageSize
	}

	return result, page
}

func GetPageNum(page int) int {
	return page * setting.AppSetting.PageSize
}
