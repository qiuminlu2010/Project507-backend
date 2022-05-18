package gin_http

import (
	"fmt"
	"net/http"
	"qiu/blog/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Bind(c *gin.Context, data interface{}) (int, int) {
	var err error
	if c.ContentType() == "application/x-www-form-urlencoded" {
		// fmt.Println("绑定form")
		err = c.Bind(data)
	} else {
		// fmt.Println("绑定json")
		err = c.BindJSON(data)
	}
	if err != nil {
		// fmt.Println("绑定数据", data)
		// fmt.Println("绑定错误", err)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	validate := validator.New()
	err = validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("%v should %v %v, but got %v\n", err.Namespace(), err.Tag(), err.Param(), err.Value())
			// logging.Info("%v should %v %v, but got %v", err.Namespace(), err.Tag(), err.Param(), err.Value())
		}
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
	// fmt.Println("绑定数据", data)
	return http.StatusOK, e.SUCCESS
}
