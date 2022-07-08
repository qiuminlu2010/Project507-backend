package v1

import (
	"fmt"
	"log"
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	chat "qiu/blog/service/chat"
	param "qiu/blog/service/param"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var wsUpgrader = websocket.Upgrader{}

// @Summary 获取文章列表
// @Produce  json
// @Param from_uid path int true "发送用户ID"
// @Param to_uid path int true "接收用户ID"
// @Param token header string true "token"
// @Router /api/v1/chat/{from_uid}/{to_uid} [get]
func ChatHandler(c *gin.Context) {

	params := param.ChatClientParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	fmt.Println("绑定参数", params)
	// uid := c.Query("from_uid") // 自己的id
	// toUid := c.Query("to_Uid") // 对方的id
	// wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	// conn, err := websocket.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("websocket建立连接错误", err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个用户实例
	client := &chat.Client{
		FromUid: params.FromUid,
		ToUid:   params.ToUid,
		Socket:  conn,
		Send:    make(chan []byte),
	}
	// fmt.Println("绑定client", client)
	// 用户注册到用户管理上
	chat.Manager.Register <- client
	go client.Read()
	go client.Write()
}
