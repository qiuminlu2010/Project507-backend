package service

import (
	"fmt"
	"qiu/backend/model"
	"qiu/backend/pkg/e"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/redis"
	"strconv"
	"strings"
)

var (
	keyPattern = "%s:%d:%s"
)

func GetModelIdKey(modelName string, modelId int) string {
	return fmt.Sprintf("%s:%d", modelName, modelId)
}

func GetModelFieldKey(modelName string, modelId uint, fieldName string) string {
	return fmt.Sprintf(keyPattern, modelName, modelId, fieldName)
}

func GetMessageKey(modelName string, modelId uint, fieldName string) string {
	return fmt.Sprintf("message:%s:%d:%s", modelName, modelId, fieldName)
}

//uid:2:like_articles
func GetArticleListParamsKey(pageNum int, pageSize int) string {
	return fmt.Sprintf("%s:page_num:%d:page_size:%d", "article_list", pageNum, pageSize)
}

//getModelsFromCache (modelName, modelIds)
//SetUserInfoCache
func FlushArticleLikeUsers() error {
	log.Logger.Info("FlushArticleLikeUsers")

	pattern := fmt.Sprintf("%s*%s", e.CACHE_MESSAGE, e.CACHE_LIKEUSERS)
	data := redis.ScanHashByPattern(pattern)

	for key, value := range data {

		// fmt.Println("ScanSetByPattern", key, value)
		var likeUsers []*model.ArticleLikeUsers
		var unlikeUsers []uint

		//message:article:2:like_users  1600000
		articleId, _ := strconv.Atoi(strings.Split(key, ":")[2])

		for k, v := range value.(map[string]string) {
			userId, _ := strconv.Atoi(k)
			v, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
			if v > 0 {
				likeUsers = append(likeUsers, &model.ArticleLikeUsers{UserID: userId, ArticleID: articleId, CreatedAt: v})
			} else if v < 0 {
				unlikeUsers = append(unlikeUsers, uint(userId))
			}
		}

		if len(likeUsers) > 0 {
			// fmt.Println("likeUserId", articleId, likeUsers)
			if err := model.AddArticleLikeUsers(likeUsers); err != nil {
				panic(err)
			}
		}
		if len(unlikeUsers) > 0 {
			// fmt.Println("unlikeUserId", articleId, unlikeUsers)
			if err := model.DeleteArticleLikeUsers(uint(articleId), unlikeUsers); err != nil {
				panic(err)
			}
		}
		redis.Del(key)
	}
	log.Logger.Info("FlushArticleLikeUsers", "OK")
	return nil
}

func FlushUserFollows() error {
	log.Logger.Info("FlushUserFollows")
	pattern := fmt.Sprintf("%s*%s", e.CACHE_MESSAGE, e.CACHE_FOLLOWS)
	data := redis.ScanHashByPattern(pattern)
	for key, value := range data {
		var follow []int
		var unFollow []int
		userId, _ := strconv.Atoi(strings.Split(key, ":")[2])
		for k, v := range value.(map[string]string) {
			followId, _ := strconv.Atoi(k)
			if v == "1" {
				follow = append(follow, followId)
			} else if v == "0" {
				unFollow = append(unFollow, followId)
			}
		}
		if len(follow) > 0 {
			if err := model.FollowUsers(uint(userId), follow); err != nil {
				panic(err)
			}
		}
		if len(unFollow) > 0 {
			if err := model.UnFollowUsers(uint(userId), unFollow); err != nil {
				panic(err)
			}
		}
		redis.Del(key)
	}
	log.Logger.Info("FlushUserFollows", "OK")
	return nil
}
