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
	apiv1 := r.Group("/api/v1")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//用户类
	r.POST("/user/register", v1.Register)
	r.POST("/user/login", v1.Login)
	r.POST("/user/refreshToken/:id", v1.RefreshToken)
	apiv1.GET("/tag/getList/", v1.GetTags)

	apiv1.GET("/article/getList/", v1.GetArticles)
	apiv1.Use(middleware.JWT())
	{
		//标签类
		apiv1.POST("/tag/add", v1.AddTag)
		apiv1.DELETE("/tag/delete/:id", v1.DeleteTag)
		apiv1.PUT("/tag/update/:id", v1.EditTag)
		apiv1.POST("/tag/recover/:id", v1.RecoverTag)

		//文章类
		apiv1.POST("/article/add", v1.AddArticle)
		apiv1.GET("/article/get/:id", v1.GetArticle)
		apiv1.POST("/article/addTags/:id", v1.AddArticleTags)
		apiv1.DELETE("/article/delete/:id", v1.DeleteArticle)
		apiv1.POST("/article/recover/:id", v1.RecoverArticle)

		//上传图片
		r.POST("/upload", v1.UploadImage)
	}
	r.Use(middleware.JWT())
	{
		//用户类
		r.DELETE("/user/delete", v1.DeleteUser)
		r.PUT("/user/update", v1.UpdatePassword)
	}
	r.Use(middleware.CORS())
	return r
}
