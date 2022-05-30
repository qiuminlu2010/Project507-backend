package v1

import (
	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/util"
	"qiu/blog/service"

	"qiu/blog/pkg/setting"

	"net/http"

	"qiu/blog/pkg/e"

	gin_http "qiu/blog/pkg/http"
)

//获取多个文章标签
// @Summary 获取标签列表
// @Produce  json
// @Param page query int false "Page"
// @Success 200 {object}  gin_http.ResponseJSON
// @Router /api/v1/tag/getList [get]
func GetTags(c *gin.Context) {
	tagService := service.GetTagService()

	tagService.PageNum = util.GetPage(c)
	tagService.PageSize = setting.AppSetting.PageSize

	tags := tagService.Get()
	gin_http.Response(c, http.StatusOK, e.SUCCESS, tags)
}

// @Summary 获取该标签的所有文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object}  gin_http.ResponseJSON
// @Router /api/v1/tag/getArticles/{id} [get]
func GetTagArticles(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	articles, err := tagService.GetArticles()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, articles)
}

// @Summary 新增标签
// @Produce  json
// @Param name formData string true "Name"
// @Param created_by formData string false "CreatedBy"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  10001 {object} gin_http.ResponseJSON
// @Failure  10006 {object} gin_http.ResponseJSON
// @Router /api/v1/tag/add [post]
func AddTag(c *gin.Context) {

	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	// err := tagService.Valid()
	// if err != nil {
	// 	gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// 	return
	// }
	state := tagService.ExistTag()
	if state {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_EXIST_TAG, nil)
		return
	}

	claims := tagService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	//TODO:还需验证用户是否存在

	// created_by := tagService.GetCreatedBy()
	// if created_by == "" {
	// 	tagService.SetCreatedBy(claims.Username)
	// } else {
	// 	if created_by != claims.Username {
	// 		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
	// 		return
	// 	}
	// }

	err := tagService.Add()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改标签
// @Produce  json
// @Param id path int true "ID"
// @Param name formData string true "Name"
// @Param modified_by formData string false "Modifiedby"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  10007 {object} gin_http.ResponseJSON
// @Router /api/v1/tag/update/{id} [put]
func EditTag(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := tagService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	// modified_by := tagService.GetModifiedBy()
	// if modified_by == "" {
	// 	tagService.SetModifiedBy(claims.Username)
	// } else {
	// 	if modified_by != claims.Username {
	// 		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
	// 		return
	// 	}
	// }

	if err := tagService.Update(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

// @Summary 删除标签
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  10008 {object} gin_http.ResponseJSON
// @Router /api/v1/tag/delete/{id} [delete]
func DeleteTag(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	// id, _ := strconv.Atoi(c.Param("id"))
	// tagService.SetId(id)

	// err := tagService.Valid()
	// if err != nil {
	// 	gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// 	return
	// }

	if err := tagService.Delete(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 恢复标签
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Router /api/v1/tag/recover/{id} [post]
func RecoverTag(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	if err := tagService.Recovery(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REC_TAG_FAIL, nil)
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 清理标签(硬删除)
// @Produce  json
// @Param id path int true "ID"
// @Param token header string true "token"
// @Router /api/v1/tag/clear/{id} [delete]
func ClearTag(c *gin.Context) {
	tagService := service.GetTagService()
	httpCode, errCode := tagService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	if err := tagService.Clear(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}
