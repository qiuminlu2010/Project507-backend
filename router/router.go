package router

import (
	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/setting"

	v1 "qiu/blog/router/api/v1"

	middleware "qiu/blog/middleware"

	_ "qiu/blog/docs"

	"github.com/swaggo/gin-swagger"

	swaggerFiles "github.com/swaggo/files"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//注册用户
	//curl "http://127.0.0.1:8000/api/v1/register?username=qxy1&password=123456" --include --header "Content-Type: application/json" --request "POST"
	r.POST("/register", v1.Register)
	//登录用户
	// r.GET("/login", v1.Login)
	r.POST("/login", v1.Login)
	apiv1 := r.Group("/api/v1")
	//获取标签列表
	apiv1.GET("/tag/getList/", v1.GetTags)

	//更新指定文章
	// apiv1.PUT("/articles/:id", v1.EditArticle)
	//删除指定文章
	// apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	apiv1.Use(middleware.JWT())
	{
		//新建标签
		apiv1.POST("/tag/add", v1.AddTag)
		//删除指定标签
		apiv1.DELETE("/tag/delete/:id", v1.DeleteTag)
		//更新指定标签
		apiv1.PUT("/tag/update/:id", v1.EditTag)

		//新建文章
		apiv1.POST("/article/add", v1.AddArticle)
		//获取文章列表
		apiv1.GET("/article/getList/", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/article/get/:id", v1.GetArticle)
		//添加文章标签
		apiv1.POST("/article/addTags/:id", v1.AddArticleTags)

		//上传图片
		r.POST("/upload", v1.UploadImage)
	}
	r.Use(middleware.JWT())
	{
		//删除用户
		r.DELETE("/delete_user", v1.DeleteUser)
		//修改用户密码
		r.PUT("/update_password", v1.UpdatePassword)
	}
	r.Use(middleware.CORS())
	return r
}
