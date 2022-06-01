package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_TAG:                 "已存在该标签名称",
	ERROR_EXIST_TAG_FAIL:            "获取已存在标签失败",
	ERROR_NOT_EXIST_TAG:             "该标签不存在",
	ERROR_GET_TAGS_FAIL:             "获取所有标签失败",
	ERROR_COUNT_TAG_FAIL:            "统计标签失败",
	ERROR_ADD_TAG_FAIL:              "新增标签失败",
	ERROR_EDIT_TAG_FAIL:             "修改标签失败",
	ERROR_DELETE_TAG_FAIL:           "删除标签失败",
	ERROR_EXPORT_TAG_FAIL:           "导出标签失败",
	ERROR_IMPORT_TAG_FAIL:           "导入标签失败",
	ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	ERROR_ADD_ARTICLE_FAIL:          "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:       "删除文章失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_LOGIN:                     "登录失败",
	ERROR_REGISTER:                  "注册失败",
	ERROR_EXIST_USER:                "用户已存在",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",
	ERROR_DELETE_USER_FAIL:          "删除用户失败",
	ERROR_UPDATE_USER_FAIL:          "更新用户信息失败",
	ERROR_REDIS_KEY_FAIL:            "缓存获取失败",
	ERROR_ADD_ARTICLE_TAG_FAIL:      "添加文章标签失败",
	ERROR_REC_ARTICLE_FAIL:          "恢复文章失败",
	ERROR_REC_TAG_FAIL:              "恢复标签失败",
	ERROR_UUID_EXPIRE:               "UUID过期",
	ERROR_GET_USERID_FAIL:           "获取文章的用户ID失败",
	ERROR_DELETE_ARTICLE_TAG_FAIL:   "删除文章标签失败",
	ERROR_IMAGE_LOST:                "上传图片缺失",
	ERROR_UPLOAD_IMAGE_FAIL:         "上传图片失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
