package v1

import (
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	log "qiu/blog/pkg/logging"
	"qiu/blog/pkg/setting"
	msg "qiu/blog/service/msg"
	param "qiu/blog/service/param"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var wsUpgrader = websocket.Upgrader{}

// @Summary 私信
// @Produce  json
// @Param from_uid path int true "发送用户ID"
// @Param to_uid path int true "接收用户ID"
// @Param token header string true "token"
// @Router /api/v1/chat/{from_uid}/{to_uid} [get]
func Chat(c *gin.Context) {

	params := param.ChatClientParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定参数", params)
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
		log.Logger.Error("websocket建立连接错误", err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个用户实例
	client := &msg.Client{
		FromUid: params.FromUid,
		ToUid:   params.ToUid,
		Socket:  conn,
		Send:    make(chan []byte),
	}
	// fmt.Println("绑定client", client)
	// 用户注册到用户管理上
	msg.Manager.Register <- client
	go client.Read()
	go client.Write()
}

// @Summary 历史消息
// @Produce  json
// @Param from_uid query int true "发送用户ID"
// @Param to_uid query int true "接收用户ID"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Param token header string true "token"
// @Router /api/v1/chat/history [get]
func GetMessage(c *gin.Context) {
	params := param.ChatMessageGetParams{}
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize
	log.Logger.Debug("绑定参数", params)
	messages, err := model.GetMessages(params.FromUid, params.ToUid, params.PageNum, params.PageSize)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	data := make(map[string]interface{})
	data["datalist"] = messages
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 聊天室
// @Produce  json
// @Param from_uid query int true "发送用户ID"
// @Param token header string true "token"
// @Router /api/v1/chat/history [get]
func ChatRoom(c *gin.Context) {

	// uid, err := strconv.Atoi(c.Query("from_uid"))
	// if err != nil {
	// 	log.Logger.Error("绑定错误", err)
	// 	gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// 	return
	// }
	// log.Logger.Debug("绑定参数", params)
	// // uid := c.Query("from_uid") // 自己的id
	// // toUid := c.Query("to_Uid") // 对方的id
	// // wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// // conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	// conn, err := (&websocket.Upgrader{
	// 	CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
	// 		return true
	// 	}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	// // conn, err := websocket.Upgrade(c.Writer, c.Request, nil)
	// if err != nil {
	// 	log.Logger.Error("websocket建立连接错误", err)
	// 	http.NotFound(c.Writer, c.Request)
	// 	return
	// }
	// // 创建一个用户实例
	// client := &msg.Client{
	// 	FromUid: params.FromUid,
	// 	ToUid:   params.ToUid,
	// 	Socket:  conn,
	// 	Send:    make(chan []byte),
	// }
	// // fmt.Println("绑定client", client)
	// // 用户注册到用户管理上
	// msg.Manager.Register <- client
	// go client.Read()
	// go client.Write()
}
