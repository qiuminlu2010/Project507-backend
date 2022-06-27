package service

import (
	"fmt"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"strconv"
	"strings"
)

var (
	keyPattern = "%s:%d:%s"
)

func GetModelKey(modelName string, modelId uint, fieldName string) string {
	return fmt.Sprintf(keyPattern, modelName, modelId, fieldName)
}

func GetMessageKey(modelName string, modelId uint, fieldName string) string {
	return fmt.Sprintf("message:%s:%d:%s", modelName, modelId, fieldName)
}

//uid:2:like_articles
func GetArticleListParamsKey(pageNum int, pageSize int) string {
	return fmt.Sprintf("%s:page_num:%d:page_size:%d", "article_list", pageNum, pageSize)
	// if a.ID > 0 {
	// 	keys = append(keys, strconv.Itoa(a.ID))
	// }
	// if a.TagID > 0 {
	// 	keys = append(keys, strconv.Itoa(a.TagID))
	// }
	// if a.State >= 0 {
	// 	keys = append(keys, strconv.Itoa(a.State))
	// }
}

func FlushArticleLikeUsers() error {
	pattern := fmt.Sprintf("%s*%s", e.CACHE_MESSAGE, e.CACHE_LIKEUSERS)
	data := redis.ScanSetByPattern(pattern)
	for key, value := range data {
		fmt.Println("ScanSetByPattern", key, value)
		var likeUsers []uint
		var unlikeUsers []uint
		articleId, _ := strconv.Atoi(strings.Split(key, ":")[2])
		// article := model.Article{}
		// article.ID = uint(articleId)
		for _, v := range value {
			userId, _ := strconv.Atoi(v)
			cache_key := GetModelKey(e.CACHE_ARTICLE, uint(userId), e.CACHE_LIKEUSERS)
			// user := model.User{}
			// user.ID = uint(userId)
			if redis.GetBit(cache_key, int64(userId)) == 1 {
				likeUsers = append(likeUsers, uint(userId))
			} else {
				unlikeUsers = append(unlikeUsers, uint(userId))
			}
		}

		if len(likeUsers) > 0 {
			fmt.Println("likeUserId", articleId, likeUsers)
			if err := model.AddArticleLikeUsers(uint(articleId), likeUsers); err != nil {
				panic(err)
			}
			// if err := model.GetArticleLikeCount(article); err != nil {
			// 	panic(err)
			// }
			// addlikeArticleUsers(article, likeUsers)
		}
		if len(unlikeUsers) > 0 {
			fmt.Println("unlikeUserId", articleId, unlikeUsers)
			if err := model.DeleteArticleLikeUsers(uint(articleId), unlikeUsers); err != nil {
				panic(err)
			}

		}
		redis.Del(key)
	}
	return nil
}

// func addlikeArticleUsers(article *Article, users []*User) {

// }

// func deletelikeArticleUsers(data []model.ArticleIdUserId) {
// 	if err := model.DeleteArticleLikeUsers(data); err != nil {
// 		panic(err)
// 	}
// }
