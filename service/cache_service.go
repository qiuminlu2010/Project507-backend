package service

import "fmt"

var (
	keyPattern = "{%s:%d}:%s"
)

func GetKeyName(modelName string, modelId uint, fieldName string) string {
	return fmt.Sprintf(keyPattern, modelName, modelId, fieldName)
}

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
	return nil
}