package gin_http

import (
	"net/http"
	"qiu/blog/pkg/e"

	"github.com/gin-gonic/gin"
)

func Bind(c *gin.Context, data interface{}) (int, int) {
	var err error
	if c.ContentType() == "application/x-www-form-urlencoded" {
		err = c.Bind(data)
	} else {
		err = c.BindJSON(data)
	}
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	// valid := validation.Validation{}
	// check, err := valid.Valid(form)
	// if err != nil {
	// 	return http.StatusInternalServerError, e.ERROR
	// }
	// if !check {
	// 	MarkErrors(valid.Errors)
	// 	return http.StatusBadRequest, e.INVALID_PARAMS
	// }
	return http.StatusOK, e.SUCCESS
}
