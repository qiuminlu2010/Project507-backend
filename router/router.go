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
	r.POST("/user/refreshToken/:id", v1.RefreshToken)
	//标签类
	apiv1.GET("/tag/getArticles/:id", v1.GetTagArticles)
	apiv1.GET("/tag/getList/", v1.GetTags)
	//文章类
	apiv1.GET("/article/list/", v1.GetArticles)
	apiv1.GET("/article/get/:id", v1.GetArticle)

	apiv1.Use(middleware.JWT())
	{
		//标签类
		apiv1.POST("/tag/add", v1.AddTag)
		apiv1.DELETE("/tag/delete/:id", v1.DeleteTag)
		apiv1.PUT("/tag/update/:id", v1.EditTag)
		apiv1.POST("/tag/recover/:id", v1.RecoverTag)
		apiv1.DELETE("/tag/clear/:id", v1.ClearTag)

		//文章类
		apiv1.POST("/article/add", v1.AddArticle)
		apiv1.POST("/article/addTags/:id", v1.AddArticleTags)
		apiv1.DELETE("/article/delete/:id", v1.DeleteArticle)
		apiv1.DELETE("/article/deleteTags/:id", v1.DeleteArticleTags)
		apiv1.POST("/article/recover/:id", v1.RecoverArticle)
		apiv1.PUT("/article/update/:id/state", v1.UpdateArticle)
		apiv1.POST("/article/like/:id", v1.LikeArticle)

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
