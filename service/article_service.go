package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	"strconv"
	"time"
)

type ArticleService struct {
	BaseService
	ArticleParams
}

func GetArticleService() *ArticleService {
	s := ArticleService{}
	s.model = &s
	return &s
}

func (s *ArticleService) GetTagName() []string {
	return s.TagName
}

func (s *ArticleService) CheckTagName() (int, int) {
	for _, tag_name := range s.TagName {
		if !model.ExistTagByName(tag_name) {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
	}
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) AddArticleTags() (int, int) {
	var tags []model.Tag
	for _, tag_id := range s.TagID {
		if !model.ExistTagByID(tag_id) {
			return http.StatusBadRequest, e.ERROR_NOT_EXIST_TAG
		}
		tag := model.Tag{}
		tag.ID = uint(tag_id)
		tags = append(tags, tag)
	}
	if err := model.AddArticleTags(s.Id, tags); err != nil {
		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
	}
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) DeleteArticleTags() (int, int) {
	var tags []model.Tag
	for _, tag_id := range s.TagID {
		tag := model.Tag{}
		tag.ID = uint(tag_id)
		tags = append(tags, tag)
	}
	if err := model.DeleteArticleTag(s.Id, tags); err != nil {
		return http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL
	}
	return http.StatusOK, e.SUCCESS
}

func (s *ArticleService) Add() error {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		if err != nil {
			return err
		}
		tag := model.Tag{}
		tag.ID = tag_id
		tags = append(tags, tag)
	}

	if err := model.AddArticle(
		model.Article{
			OwnerID: s.UserID,
			Content: s.Content,
		}, tags); err != nil {
		return err
	}
	return nil
}

func (s *ArticleService) AddArticleWithImg() error {
	var tags []model.Tag
	for _, tag_name := range s.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		tag := model.Tag{}
		if err != nil {
			tag.Name = tag_name
		} else {
			tag.ID = tag_id
		}
		tags = append(tags, tag)
	}
	var imgs []model.Image
	for _, img_name := range s.ImgName {
		imgs = append(imgs, model.Image{Filename: img_name})
	}
	if err := model.AddArticleWithImg(
		model.Article{
			OwnerID: s.UserID,
			Content: s.Content,
			Title:   s.Title,
		}, tags, imgs); err != nil {
		return err
	}
	return nil
}
func (s *ArticleService) ExistByID() bool {
	return model.ExistArticleByID(s.Id)
}

// article:id (bytes) => ArticleInfo
func getArticleCache(articleId int) (*model.ArticleInfo, error) {
	key := GetModelIdKey(e.CACHE_ARTICLE, articleId)
	var article model.ArticleInfo
	if redis.Exists(key) == 0 {
		if err := setArticleCache(articleId); err != nil {
			return &article, err
		}
	}
	if err := json.Unmarshal(redis.GetBytes(key), &article); err != nil {
		return &article, err
	}
	return &article, nil
}

// ArticleInfo => article:id (bytes)
func setArticleCache(articleId int) error {
	articleInfo, err := model.GetArticle(articleId)
	if err != nil {
		return err
	}
	key := GetModelIdKey(e.CACHE_ARTICLE, articleId)
	redis.SetBytes(key, articleInfo, 6*time.Hour)
	return nil
}

// ZSet: articles (score:timestamp member:article_id)
func setArticleIdsCache() error {
	articles, err := model.GetArticlesCache()
	if err != nil {
		panic(err)
		// return err
	}
	for _, article := range articles {
		redis.ZAdd(e.CACHE_ARTICLES, float64(article.CreatedOn), article.ID)
	}
	return nil
}

func setArticlesCache(articleIds []int) {
	articles, err := model.GetArticlesByIds(articleIds)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		redis.SetBytes(GetModelIdKey(e.CACHE_ARTICLE, int(article.ID)), article, 3*time.Hour)
	}
}

// func getArticlesCache(stringIds []string) ([]*model.ArticleInfo, error) {
// 	// var articles []*model.ArticleInfo
// 	var ids []int
// 	articleIds := util.StringsToInts(stringIds)
// 	for _, articleId := range articleIds {
// 		// articleId, _ := strconv.Atoi(v)
// 		// articleInfo, err := getArticleCache(articleId)
// 		// if err != nil {
// 		// 	return articles, err
// 		// }
// 		// if redis.Exists()
// 		ids = append(ids, articleId)
// 		// articles = append(articles, articleInfo)
// 	}
// 	articles, err := model.GetArticlesByIds(ids)
// 	if err != nil {
// 		return articles, err
// 	}
// 	return articles, nil
// }

//(article_id1, article_id2, ...) => []ArticleInfo
func getArticleInfoFromCache(stringIds []string) ([]*model.ArticleInfo, error) {
	var articles []*model.ArticleInfo
	articleIds, err := util.StringsToInts(stringIds)
	if err != nil {
		return nil, err
	}
	var setCacheIds []int
	for _, articleId := range articleIds {
		if redis.Exists(GetModelIdKey(e.CACHE_ARTICLE, articleId)) == 0 {
			setCacheIds = append(setCacheIds, articleId)
		}
	}
	if len(setCacheIds) > 0 {
		setArticlesCache(setCacheIds)
	}

	for _, articleId := range articleIds {
		articleInfo, err := getArticleCache(articleId)
		if err != nil {
			return articles, err
		}
		articles = append(articles, articleInfo)
	}
	return articles, nil
}

func (s *ArticleService) GetArticles(params ArticleGetParams) ([]*model.ArticleInfo, error) {
	var (
		err      error
		articles []*model.ArticleInfo
		// articles []*model.Article
	)
	if redis.Exists(e.CACHE_ARTICLES) == 0 {
		if err := setArticleIdsCache(); err != nil {
			return nil, err
		}
	}
	articleIds := redis.ZRevRange(e.CACHE_ARTICLES, int64(params.PageNum), int64(params.PageSize-1))
	articles, err = getArticleInfoFromCache(articleIds)
	if err != nil {
		return nil, err
	}
	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}
	return articles, nil

	// articles, err = model.GetArticles(params.PageNum, params.PageSize, nil)

	// if err != nil {
	// 	return nil, err
	// }
	// // //TODO: 写数据库 缓存一致性问题
	// // redis.SetBytes(key, articles, time.Minute*3)
	// if err = getArticleLikeInfo(articles, params.Uid); err != nil {
	// 	return nil, err
	// }

	// return articles, nil
}

//根据uid获取articles的IsLike字段,需要优化
func getArticleLikeInfo(articles []*model.ArticleInfo, uid int) error {
	for _, article := range articles {
		key := GetModelFieldKey(e.CACHE_ARTICLE, article.ID, e.CACHE_LIKEUSERS)
		if redis.Exists(key) == 0 {
			fmt.Println("SET CACHE_KEY", key)
			likeUsers, err := model.GetArticleLikeUsers(article.ID)
			if err != nil {
				return err
			}
			redis.SetBit(key, 0, 0)
			for _, user := range likeUsers {
				redis.SetBit(key, int64(user.UserId), 1)
			}
		}
		article.LikeCount = redis.BitCount(key)
		// cnt := model.GetArticleLikeCount(article)
		// fmt.Println("LikeCount", cnt)
		if uid != 0 {
			article.IsLike = redis.GetBit(key, int64(uid)) == 1
		}
	}
	return nil
}

func (s *ArticleService) UpdateArticleLike(param ArticleLikeParams) error {
	key := GetModelFieldKey(e.CACHE_ARTICLE, uint(param.Id), e.CACHE_LIKEUSERS)
	messageKey := GetMessageKey(e.CACHE_ARTICLE, uint(param.Id), e.CACHE_LIKEUSERS)
	if redis.Exists(key) == 0 {
		if err := setArticleLikeCache(uint(param.Id)); err != nil {
			return err
		}
	}

	redis.SetBit(key, int64(param.UserID), param.Type)
	m := make(map[string]interface{})
	if param.Type == 1 {
		m[strconv.Itoa(param.UserID)] = time.Now().Unix()
	} else {
		m[strconv.Itoa(param.UserID)] = -time.Now().Unix()
	}

	redis.HashSet(messageKey, m)
	// redis.SAdd(messageKey, param.UserID)
	return nil
	// return model.AddArticleLikeUser(uint(param.Id), user)
}

func setArticleLikeCache(articleId uint) error {
	key := GetModelFieldKey(e.CACHE_ARTICLE, articleId, e.CACHE_LIKEUSERS)
	likeUsers, err := model.GetArticleLikeUsers(articleId)
	if err != nil {
		return err
	}
	for _, user := range likeUsers {
		redis.SetBit(key, int64(user.UserId), 1)
	}
	return nil
}

//(userId, pageNum, pageSize) => []ArticleInfo
func (s *UserService) GetLikeArticles(params *ArticleGetParams) ([]*model.ArticleInfo, error) {

	key := GetModelFieldKey(e.CACHE_USER, uint(params.Uid), e.CACHE_LIKEARTICLES)
	var articles []*model.ArticleInfo
	var err error
	if redis.Exists(key) == 0 {
		if err = setLikeArticleCache(params.Uid); err != nil {
			return nil, err
		}
	}
	//有序集合，根据时间戳降排
	articleIds := redis.ZRevRange(key, int64(params.PageNum), int64(params.PageSize-1))
	articles, err = getArticleInfoFromCache(articleIds)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

//设置缓存： user:id:like_articles (score:time member:articleId)
func setLikeArticleCache(userId int) error {
	key := GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_LIKEARTICLES)
	likes, err := model.GetLikeArticles(uint(userId))
	if err != nil {
		return err
	}
	for _, like := range likes {
		redis.ZAdd(key, float64(like.CreatedAt), like.ArticleID)
	}
	return nil
}

func (s *ArticleService) Delete() error {
	return model.DeleteArticle(s.Id)
}

func (s *ArticleService) Update() error {
	data := make(map[string]interface{})
	data["state"] = s.State
	data["content"] = s.Content
	return model.UpdateArticle(s.Id, data)
}
func (s *ArticleService) Count(data map[string]interface{}) int64 {
	if redis.Exists(e.CACHE_ARTICLES) == 0 {
		if err := setArticleIdsCache(); err != nil {
			return 0
		}
	}
	return redis.ZCard(e.CACHE_ARTICLES)
	// return model.GetArticleTotal(data)
}

// func (s *ArticleService) Clear() error {
// 	return model.CleanAllArticle()
// }

func (s *ArticleService) Recovery() error {
	return model.RecoverArticle(s.Id)
}

func (s *ArticleService) GetUserID() (uint, error) {
	return model.GetArticleUserID(s.Id)
}

// func (a *ArticleService) GetArticleKey() string {
// 	return e.CACHE_ARTICLE + "_" + strconv.Itoa(a.ID)
// }
