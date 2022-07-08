package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	articleService "qiu/blog/service/article"
	param "qiu/blog/service/param"
	service "qiu/blog/service/tag"

	"qiu/blog/pkg/setting"

	"net/http"

	"qiu/blog/pkg/e"

	gin_http "qiu/blog/pkg/http"
)

//获取多个文章标签
// @Summary 获取标签列表
// @Produce  json
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Success 200 {object}  gin_http.ResponseJSON
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	tagService := service.GetTagService()
	params := param.PageGetParams{}
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	// tagService.PageNum, _ = util.GetPage(c)
	// tagService.PageSize = setting.AppSetting.PageSize
	tags := tagService.Get(&params)
	gin_http.Response(c, http.StatusOK, e.SUCCESS, tags)
}

// @Summary 获取该标签的所有文章
// @Produce  json
// @Param tag_name query string true "tag_name"
// @Param uid query int false "user_id"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Success 200 {object}  gin_http.ResponseJSON
// @Router /api/v1/tag [get]
func GetTagArticles(c *gin.Context) {
	tagService := articleService.GetArticleService()
	params := param.TagArticleGetParams{}
	// httpCode, errCode := tagService.Bind(c)
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize
	fmt.Println("绑定数据", params)

	articles, err := tagService.GetArticlesByMultiTags(&params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["datalist"] = articles
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  10008 {object} gin_http.ResponseJSON
// @Router /api/v1/tag/{id} [delete]
func DeleteTag(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagId <= 0 {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := tagService.Delete(tagId); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 恢复标签
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Router /api/v1/tag/{id}/recover [put]
func RecoverTag(c *gin.Context) {
	tagService := service.GetTagService()
	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagId <= 0 {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := tagService.Recovery(tagId); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REC_TAG_FAIL, nil)
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 清理标签(硬删除)
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Router /api/v1/tag/{id}/clear [delete]
func ClearTag(c *gin.Context) {
	tagService := service.GetTagService()
	tagId, err := strconv.Atoi(c.Param("id"))
	if err != nil || tagId <= 0 {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if err := tagService.Clear(tagId); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 标签补全
// @Produce  json
// @Param tag_name query string true "tag_name"
// @Param page_size query int false "page_size"
// @Success 200 {object}  gin_http.ResponseJSON
// @Router /api/v1/search/tag [get]
func GetTagsByPrefix(c *gin.Context) {
	tagService := service.GetTagService()
	params := param.TagsGetParams{}
	// httpCode, errCode := tagService.Bind(c)
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	// page := params.PageNum
	// params.PageNum = params.PageNum * params.PageSize
	fmt.Println("绑定数据", params)

	tags, _ := tagService.Hint(&params)
	// tags := tagService.HintByCache(&params)
	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
	// 	return
	// }
	// data := make(map[string]interface{})
	// data["datalist"] = articles
	// data["total"] = total
	gin_http.Response(c, http.StatusOK, e.SUCCESS, tags)
}
