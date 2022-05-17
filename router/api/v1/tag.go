package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/unknwon/com"

	"qiu/blog/util"

	"qiu/blog/pkg/setting"

	"net/http"

	"qiu/blog/pkg/e"

	"qiu/blog/model"

	"github.com/astaxie/beego/validation"
)

//获取多个文章标签
// @Summary 获取多个文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = model.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = model.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//新增文章标签
//curl "http://127.0.0.1:8000/api/v1/tags?name=game&state=1&created_by=qiu" --include --header "Content-Type: application/json" --request "POST"
// @Summary 新增文章标签
// @Produce  json
// @Param name query string true "Name"
// @Param created_by query string false "CreatedBy"
// @Param token query string true "token"
// @Success 200 {string} string "ok"
// @Failure  10001 {string} string "{"code":10001,"data":{},"msg":"已存在该标签"}"
// @Failure  10000 {string} string "{"code":10001,"data":{},"msg":"TOKEN为空"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	//state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}

	valid.Required(name, "name").Message("名称不能为空")

	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	//valid.Required(createdBy, "created_by").Message("创建人不能为空")

	//valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")

	//valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		if !model.ExistTagByName(name) {
			code = e.SUCCESS
			model.AddTag(name, 1, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if model.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			model.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @Summary 删除文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param token query string true "token"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure  10000 {string} string "{"code":10001,"data":{},"msg":"TOKEN为空"}"
// @Failure  10002 {string} string "{"code":10002,"data":{},"msg":"标签不存在"}"
// @Router /api/v1/tags/:id [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if model.ExistTagByID(id) {
			model.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}
