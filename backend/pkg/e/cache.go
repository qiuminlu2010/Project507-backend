package e

import "time"

const (
	CACHE_TOKEN              = "token"
	CACHE_ARTICLE            = "article"
	CACHE_ARTICLES           = "articles"
	CACHE_TAG                = "tag"
	CACHE_TAGS               = "tags"
	CACHE_LIKEUSERS          = "like_users"
	CACHE_USER               = "user"
	CACHE_LIKEARTICLES       = "like_articles"
	CACHE_MESSAGE            = "message"
	CACHE_FOLLOWS            = "follows"
	CACHE_SESSIONS           = "sessions"
	CACHE_UNREAD_MSG         = "unread_msg"
	CHANNEL_LIKEARTICLE      = "likeArticle"
	CHANNEL_LIKECOMMENT      = "likeComment"
	DURATION_USER_INFO       = time.Hour * 3
	DURATION_ARTICLE_INFO    = time.Hour * 3
	DURATION_LIKEUSERS       = time.Hour * 3
	DURATION_FOLLOWS         = time.Hour * 3
	DURATION_LIKEARTICLES    = time.Hour * 3
	DURATION_USERARTICLES    = time.Hour * 3
	DURATION_USER_TOKEN      = time.Hour * 24
	DURATION_USER_SESSIONS   = time.Hour * 24
	DURATION_USER_UNREAD_MSG = time.Hour * 24 * 30
)
