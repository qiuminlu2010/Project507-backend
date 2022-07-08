package v1

import (
	"fmt"
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	chat "qiu/blog/service/chat"
	param "qiu/blog/service/param"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{}

// apt/v1/chat/{from}/{to}
func WsHandler(c *gin.Context) {

	params := param.ChatClientParams{}
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	// uid := c.Query("from_uid") // 自己的id
	// toUid := c.Query("to_Uid") // 对方的id
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	// conn, err := websocket.Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个用户实例
	client := &chat.Client{
		ID:     params.FromUid,
		SendID: params.ToUid,
		Socket: conn,
		Send:   make(chan []byte),
	}
	// 用户注册到用户管理上
	chat.Manager.Register <- client
	// go client.Read()
	// go client.Write()
}
