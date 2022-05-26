package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/setting"
	"qiu/blog/pkg/util"

	gin_http "qiu/blog/pkg/http"

	service "qiu/blog/service"
)

//获取单个文章
func GetArticle(c *gin.Context) {

	// articleService := service.ArticleService{}
	// httpCode, errCode := gin_http.Bind(c, &articleService)

	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }

	// exists, err := articleService.ExistByID()

	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
	// 	return
	// }
	// if !exists {
	// 	gin_http.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	// 	return
	// }
	// getArticle, err := articleService.Get()
	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
	// 	return
	// }

	// gin_http.Response(c, http.StatusOK, e.SUCCESS, getArticle)
}

//获取多个文章
func GetArticles(c *gin.Context) {

	articleService := service.GetArticleService()
	articleService.PageNum = util.GetPage(c)
	articleService.PageSize = setting.AppSetting.PageSize

	data := make(map[string]interface{})
	data["delete_on"] = 0
	data["state"] = 1
	total, err := articleService.Count(data)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetArticles(data)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data["lists"] = articles
	data["total"] = total

	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

//新增文章
func AddArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	//TODO:还需验证用户是否存在

	created_by := articleService.GetCreatedBy()
	if created_by == "" {
		articleService.SetCreatedBy(claims.Username)
	} else {
		if created_by != claims.Username {
			gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
			return
		}
	}

	err := articleService.Add()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	// tagService := tag_service.Tag{ID: form.TagID}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

func AddArticleTags(c *gin.Context) {
	fmt.Println("添加文章标签", c.Params)
	articleService := service.GetArticleService()

	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	// httpCode, errCode = articleService.Bind(c)
	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }
	if httpCode, errCode = articleService.AddTags(); errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}
