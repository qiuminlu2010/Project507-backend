package article_service

import (
	"encoding/json"
	"fmt"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/upload"
	"qiu/blog/pkg/util"
	base "qiu/blog/service/base"
	cache "qiu/blog/service/cache"
	param "qiu/blog/service/param"
	user "qiu/blog/service/user"
	"strconv"
	"strings"
	"time"
)

type ArticleService struct {
	base.BaseService
}

var articleService ArticleService

func GetArticleService() *ArticleService {
	return &articleService
}

//添加文章： 1.更新数据库 2.添加文章目录缓存 3.删除用户缓存 or 添加用户缓存
func (s *ArticleService) Add(params *param.ArticleAddParams) error {
	var tags []model.Tag
	for _, tag_name := range params.TagName {
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
	for _, img_url := range params.ImgUrl {
		imgs = append(imgs, model.Image{Url: upload.GetImagePath() + img_url, ThumbUrl: upload.GetThumbPath() + img_url})
	}

	articleId, err := model.AddArticleWithImg(
		model.Article{
			OwnerID: params.UserID,
			Content: params.Content,
			Title:   params.Title,
		}, tags, imgs)
	if err != nil {
		return err
	}
	duration := float64(time.Now().Unix())
	redis.ZAdd(e.CACHE_ARTICLES, duration, articleId)

	userKey := cache.GetModelFieldKey(e.CACHE_USER, params.UserID, e.CACHE_ARTICLES)
	// redis.Del(userKey)
	if redis.Exists(userKey) != 0 {
		redis.ZAdd(userKey, duration, articleId)
	}
	return nil
}

func (s *ArticleService) Delete(articleId int) error {
	if err := model.DeleteArticle(uint(articleId)); err != nil {
		return err
	}
	redis.Del(cache.GetModelIdKey(e.CACHE_ARTICLE, articleId))
	return nil
}

func (s *ArticleService) Update(params *param.ArticleUpdateParams) error {
	data := make(map[string]interface{})
	if params.Title != "" {
		data["title"] = params.Title
	}
	if params.Content != "" {
		data["content"] = params.Content
	}
	if params.Content != "" || params.Title != "" {
		return model.UpdateArticle(params.ArticleId, data)
	}
	return nil
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

func (s *ArticleService) Recovery(articleId int) error {
	return model.RecoverArticle(articleId)
}

func (s *ArticleService) GetUserID(articleId int) (uint, error) {
	return model.GetArticleUserID(uint(articleId))
}

func (s *ArticleService) AddTags(params *param.ArticleAddTagsParams) error {
	var tags []model.Tag
	for _, tag_name := range params.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		tag := model.Tag{}
		if err != nil {
			tag.Name = tag_name
		} else {
			tag.ID = tag_id
		}
		tags = append(tags, tag)
	}
	redis.Del(cache.GetModelIdKey(e.CACHE_ARTICLE, params.ArticleId))
	return model.AddArticleTags(params.ArticleId, tags)
}

func (s *ArticleService) DeleteTags(params *param.ArticleAddTagsParams) error {
	var tags []model.Tag
	for _, tag_name := range params.TagName {
		tag_id, err := model.GetTagIdByName(tag_name)
		tag := model.Tag{}
		if err != nil {
			tag.Name = tag_name
		} else {
			tag.ID = tag_id
		}
		tags = append(tags, tag)
	}
	redis.Del(cache.GetModelIdKey(e.CACHE_ARTICLE, params.ArticleId))
	return model.DeleteArticleTag(params.ArticleId, tags)
}

//文章列表 1.设置缓存 articles (score:timestamp memeber:article_id) 2. 设置缓存 article:id (ArticleInfo) 3.获取 []article_id 4.获取 []ArticleInfo
func (s *ArticleService) GetArticles(params *param.ArticleGetParams) ([]*model.ArticleInfo, error) {
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
	var articleIds []string
	if params.Type == 1 {
		articleIds = redis.ZRandMember(e.CACHE_ARTICLES, params.PageSize-1)
	} else {
		articleIds = redis.ZRevRange(e.CACHE_ARTICLES, int64(params.PageNum), int64(params.PageSize-1))
	}

	articles, err = getArticlesCache(articleIds, e.CACHE_ARTICLES)
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

//点赞操作： 1.修改缓存 article:id:like_users  2.修改缓存 user:id:like_articles 3 添加消息缓存 message:article:id:like_users
func (s *ArticleService) UpdateArticleLike(param *param.ArticleLikeParams) error {
	acticleKey := cache.GetModelFieldKey(e.CACHE_ARTICLE, uint(param.ArticleId), e.CACHE_LIKEUSERS)
	messageKey := cache.GetMessageKey(e.CACHE_ARTICLE, uint(param.ArticleId), e.CACHE_LIKEUSERS)
	userKey := cache.GetModelFieldKey(e.CACHE_USER, uint(param.UserId), e.CACHE_LIKEARTICLES)
	if redis.Exists(acticleKey) == 0 {
		if err := setArticleLikeCache(uint(param.ArticleId)); err != nil {
			return err
		}
	}
	if redis.Exists(userKey) == 0 {
		if err := setUserLikeArticleCache(param.UserId); err != nil {
			return err
		}
	}

	redis.SetBit(acticleKey, int64(param.UserId), param.Type)

	m := make(map[string]interface{})
	if param.Type == 1 {
		m[strconv.Itoa(param.UserId)] = time.Now().Unix()
		redis.ZAdd(userKey, float64(time.Now().Unix()), param.ArticleId)
	} else {
		m[strconv.Itoa(param.UserId)] = -time.Now().Unix()
		redis.ZRem(userKey, param.ArticleId)
	}

	redis.HashSet(messageKey, m)
	// redis.SAdd(messageKey, param.UserID)
	return nil
	// return model.AddArticleLikeUser(uint(param.Id), user)
}

//(userId, pageNum, pageSize) => []ArticleInfo
func (s *ArticleService) GetUserLikeArticles(params *param.ArticleGetParams) ([]*model.ArticleInfo, error) {
	//user:id:like_articles
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(params.Uid), e.CACHE_LIKEARTICLES)
	var articles []*model.ArticleInfo
	var err error
	if redis.Exists(key) == 0 {
		if err = setUserLikeArticleCache(params.Uid); err != nil {
			return nil, err
		}
	}
	//有序集合，根据时间戳降排
	articleIds := redis.ZRevRange(key, int64(params.PageNum), int64(params.PageSize-1))
	articles, err = getArticlesCache(articleIds, key)
	if err != nil {
		return nil, err
	}
	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}
	return articles, nil
}

func (s *ArticleService) GetArticlesByTag(params *param.TagArticleGetParams) ([]*model.ArticleInfo, error) {

	articleIds, err := model.GetTagArticleIds(params.TagName)
	if err != nil {
		return nil, err
	}

	articles, err := getArticlesCache(articleIds, e.CACHE_ARTICLES)
	if err != nil {
		return nil, err
	}

	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *ArticleService) GetArticlesByMultiTags(params *param.TagArticleGetParams) ([]*model.ArticleInfo, error) {

	tagNames := strings.Split(params.TagName, " ")
	var articleIds [][]int

	for _, tagName := range tagNames {
		articleId, err := model.GetTagArticleIds(tagName)
		if err != nil {
			return nil, err
		}
		articleIds = append(articleIds, articleId)
	}

	ans := util.Intersection(articleIds)
	articles, err := getArticlesCache(ans, e.CACHE_ARTICLES)

	if err != nil {
		return nil, err
	}

	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}

	return articles, nil
}

//(userId, pageNum, pageSize) => []ArticleInfo
func (s *ArticleService) GetUserArticles(params *param.ArticleGetParams) ([]*model.ArticleInfo, error) {
	//user:id:articles
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(params.Uid), e.CACHE_ARTICLES)
	var articles []*model.ArticleInfo
	var err error
	if redis.Exists(key) == 0 {
		if err = setUserArticleCache(params.Uid); err != nil {
			return nil, err
		}
	}
	//有序集合，根据时间戳降排
	articleIds := redis.ZRevRange(key, int64(params.PageNum), int64(params.PageSize-1))
	articles, err = getArticlesCache(articleIds, key)
	if err != nil {
		return nil, err
	}
	if err = getArticleLikeInfo(articles, params.Uid); err != nil {
		return nil, err
	}
	return articles, nil
}

//设置缓存： user:id:articles (score:time member:articleId)
func setUserArticleCache(userId int) error {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_ARTICLES)
	articles, err := model.GetUserArticles(uint(userId))
	if err != nil {
		return err
	}
	for _, article := range articles {
		redis.ZAdd(key, float64(article.CreatedOn), article.ID)
	}
	redis.Expire(key, e.DURATION_USERARTICLES)
	return nil
}

// article:id (bytes) => ArticleInfo
func getArticleCache(articleId int) (*model.ArticleInfo, error) {
	key := cache.GetModelIdKey(e.CACHE_ARTICLE, articleId)
	var article model.ArticleInfo
	if redis.Exists(key) == 0 {
		if err := setArticleCache(articleId); err != nil {
			return nil, err
		}
	}

	if redis.Exists(key) == 0 {
		return nil, nil
	}
	redis.Expire(key, e.DURATION_ARTICLE_INFO)
	if err := json.Unmarshal(redis.GetBytes(key), &article); err != nil {
		return nil, err
	}
	return &article, nil
}

// ArticleInfo => article:id (bytes)
func setArticleCache(articleId int) error {
	articleInfo, err := model.GetArticle(articleId)
	if err != nil {
		return err
	}
	if articleInfo != nil {
		key := cache.GetModelIdKey(e.CACHE_ARTICLE, articleId)
		redis.SetBytes(key, articleInfo, e.DURATION_ARTICLE_INFO)
	}
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

//(article_id, ...) => (article:id ArticleInfo)
func setArticlesCache(articleIds []int) {
	articles, err := model.GetArticlesByIds(articleIds)
	if err != nil {
		panic(err)
	}
	for _, article := range articles {
		redis.SetBytes(cache.GetModelIdKey(e.CACHE_ARTICLE, int(article.ID)), article, e.DURATION_ARTICLE_INFO)
	}
}

//(article_id1, article_id2, ...) => []ArticleInfo
func getArticlesCache(ids interface{}, cachekey string) ([]*model.ArticleInfo, error) {
	var articles []*model.ArticleInfo
	var articleIds []int
	var err error
	switch ids := ids.(type) {
	case []string:
		articleIds, err = util.StringsToInts(ids)
		if err != nil {
			return nil, err
		}
	case []int:
		articleIds = ids
	}

	for _, articleId := range articleIds {
		articleInfo, err := getArticleCache(articleId)
		if err != nil {
			return articles, err
		}
		if articleInfo == nil {
			redis.ZRem(cachekey, articleId)
		} else {
			articles = append(articles, articleInfo)
		}
	}

	return articles, nil
}

//根据uid获取articles的IsLike字段,需要优化
func getArticleLikeInfo(articles []*model.ArticleInfo, uid int) error {
	for _, article := range articles {
		key := cache.GetModelFieldKey(e.CACHE_ARTICLE, article.ID, e.CACHE_LIKEUSERS)
		if redis.Exists(key) == 0 {
			fmt.Println("SET CACHE_KEY", key)
			if err := setArticleLikeCache(article.ID); err != nil {
				return err
			}
			// likeUsers, err := model.GetArticleLikeUsers(article.ID)
			// if err != nil {
			// 	return err
			// }
			// redis.SetBit(key, 0, 0)
			// for _, user := range likeUsers {
			// 	redis.SetBit(key, int64(user.UserId), 1)
			// }
		}

		article.LikeCount = redis.BitCount(key)
		// cnt := model.GetArticleLikeCount(article)
		// fmt.Println("LikeCount", cnt)
		if uid != 0 {
			article.IsLike = redis.GetBit(key, int64(uid)) == 1
		}
		userInfo := user.GetUserCache(int(article.OwnerID))
		article.OwnerName = userInfo.Name
		article.OwnerUsername = userInfo.Username
		article.OwnerAvatar = userInfo.Avatar
	}
	return nil
}

// 1.获取[]like_id 2.设置位图
func setArticleLikeCache(articleId uint) error {
	key := cache.GetModelFieldKey(e.CACHE_ARTICLE, articleId, e.CACHE_LIKEUSERS)
	likeUsers, err := model.GetArticleLikeUsers(articleId)
	if err != nil {
		return err
	}
	redis.SetBit(key, 0, 0)
	for _, user := range likeUsers {
		redis.SetBit(key, int64(user.UserId), 1)
	}
	redis.Expire(key, e.DURATION_LIKEUSERS)
	return nil
}

//设置缓存： user:id:like_articles (score:time member:articleId)
func setUserLikeArticleCache(userId int) error {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(userId), e.CACHE_LIKEARTICLES)
	likes, err := model.GetLikeArticles(uint(userId))
	if err != nil {
		return err
	}
	for _, like := range likes {
		redis.ZAdd(key, float64(like.CreatedAt), like.ArticleID)
	}
	redis.Expire(key, e.DURATION_LIKEARTICLES)
	return nil
}
