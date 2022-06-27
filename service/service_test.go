package service

import (
	"reflect"
	"testing"

	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/setting"

	"qiu/blog/model"
	"qiu/blog/pkg/logging"
)

func TestFlushArticleLikeUsers(t *testing.T) { // 测试函数名必须以Test开头，必须接收一个*testing.T类型参数
	logging.Setup()
	setting.Setup()
	model.Setup()
	redis.Setup()
	if !reflect.DeepEqual(nil, FlushArticleLikeUsers()) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Error("TestFlushArticleLikeUsers:", "Failed") // 测试失败输出错误提示
	}
}
