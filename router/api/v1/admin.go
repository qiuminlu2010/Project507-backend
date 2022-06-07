package v1

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"

	gin_http "qiu/blog/pkg/http"
)

// @Summary 获取后台管理菜单
// @Produce  json
// @Param token header string true "token"
// @Router /admin/menu/list [get]
func GetAdminMenu(c *gin.Context) {
	jsonFile, err := os.Open("conf/admin_menu_list.json")

	// 最好要处理以下错误
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADMIN_MENU_LIST_FAIL, nil)
		return
	}
	defer jsonFile.Close()

	data := make(map[string]interface{})
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&data)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADMIN_MENU_LIST_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data["data"])
}
