package router

import (
	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/setting"

	v1 "qiu/blog/api/v1"

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
	r.GET("/data/:imgType/:imgName", v1.DownloadImg)
	//用户类
	r.POST("/user/register", v1.Register)
	r.POST("/user/login", v1.Login)
	r.POST("/user/:id/refreshToken", v1.RefreshToken)
	//标签类
	apiv1.GET("/tag/:id/articles", v1.GetTagArticles)
	apiv1.GET("/tag", v1.GetTags)
	//文章类
	apiv1.GET("/article", v1.GetArticles)
	apiv1.GET("/article/:id", v1.GetArticle)

	apiv1.Use(middleware.JWT())
	{
		//标签类
		apiv1.POST("/tag", v1.AddTag)
		apiv1.DELETE("/tag/:id", v1.DeleteTag)
		apiv1.PUT("/tag/:id", v1.EditTag)
		apiv1.POST("/tag/:id/recover", v1.RecoverTag)
		apiv1.DELETE("/tag/:id/clear", v1.ClearTag)

		//文章类
		apiv1.POST("/article", v1.AddArticle)
		apiv1.POST("/article/:id/addTags", v1.AddArticleTags)
		apiv1.DELETE("/article/:id", v1.DeleteArticle)
		apiv1.DELETE("/article/:id/deleteTags", v1.DeleteArticleTags)
		apiv1.POST("/article/:id/recover", v1.RecoverArticle)
		apiv1.PUT("/article/:id/state", v1.UpdateArticle)
		apiv1.POST("/article/:id/like", v1.LikeArticle)

		//用户类
		apiv1.POST("/user/:id/follow", v1.FollowUser)
		//上传图片
		r.POST("/upload", v1.UploadImage)
	}
	r.Use(middleware.JWT())
	{
		//用户类
		r.DELETE("/user/delete/:id", v1.DeleteUser)
		r.PUT("/user/update/:id/password", v1.UpdatePassword)
		r.PUT("/user/update/:id/state", v1.UpdateUserState)
		r.GET("/user/list", v1.GetUserList)

		//后台管理
		r.GET("/admin/menu/list", v1.GetAdminMenu)
	}
	r.Use(middleware.CORS())
	return r
}
