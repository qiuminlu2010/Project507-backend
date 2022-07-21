package v1

import (
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	log "qiu/blog/pkg/logging"
	"qiu/blog/pkg/setting"
	msg "qiu/blog/service/msg"
	param "qiu/blog/service/param"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var wsUpgrader = websocket.Upgrader{}

// @Summary 私信
// @Produce  json
// @Param id path int true "发送用户ID"
// @Param token header string true "token"
// @Router /api/v1/msg/{id}/chat [get]
func Chat(c *gin.Context) {

	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil || uid <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定参数", uid)
	token := c.Param("token")
	if !msg.CheckToken(token, uid) {
		log.Logger.Error("绑定错误", token)
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
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
		Uid:    uid,
		Socket: conn,
		Send:   make(chan []byte),
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
// @Param offset query int false "offset"
// @Param limit query int false "limit"
// @Param token header string true "token"
// @Router /api/v1/msg/history [get]
func GetMessage(c *gin.Context) {
	params := param.MessageGetParams{}
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if params.Limit == 0 {
		params.Limit = setting.AppSetting.PageSize
	}
	log.Logger.Debug("绑定参数", params)
	messages, err := msg.GetMessages(&params)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	}
	// userInfo := user.GetUserCache(params.ToUid)
	data := make(map[string]interface{})
	// messagesInfo := make(map[string]interface{})
	// messagesInfo["userInfo"] = userInfo
	// messagesInfo["messages"] = messages
	data["datalist"] = messages
	// data["total"] = total
	data["offset"] = params.Offset
	data["limit"] = params.Limit
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 消息会话列表
// @Produce  json
// @Param uid query int true "用户ID"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Param token header string true "token"
// @Router /api/v1/msg/session [get]
func GetMessageSession(c *gin.Context) {
	params := param.SessionGetParams{}
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
	sessions, err := msg.GetSession(params.Uid, params.PageNum, params.PageSize)
	if err != nil {
		log.Logger.Error(e.ERROR_SESSION_GET_FAIL, err)
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_SESSION_GET_FAIL, nil)
		return
	}
	// sessions, err := model.GetSession(params.Uid, params.PageNum, params.PageSize)
	// if err != nil {
	// 	gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// }
	data := make(map[string]interface{})
	data["datalist"] = sessions
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 已读消息
// @Produce  json
// @Param uid formData int true "用户ID"
// @Param session_id formData int true "会话ID"
// @Param token header string true "token"
// @Router /api/v1/msg/read [post]
func ReadMessage(c *gin.Context) {
	params := param.UpdateUnReadMessageParams{}
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定参数", params)
	// sessions, err := model.GetSession(params.Uid, params.PageNum, params.PageSize)
	// if err != nil {
	// 	gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// }
	msg.UpdateUnReadMessage(params.Uid, params.SessionId)
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 聊天室
// @Produce  json
// @Param id path int true "发送用户ID"
// @Param token header string true "token"
// @Router /api/v1/msg/{id}/chatroom [post]
// func ChatRoom(c *gin.Context) {
// 	uid, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || uid <= 0 {
// 		log.Logger.Error("绑定错误", err)
// 		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
// 		return
// 	}
// 	log.Logger.Debug("绑定参数", uid)
// 	// uid := c.Query("from_uid") // 自己的id
// 	// toUid := c.Query("to_Uid") // 对方的id
// 	// wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	// conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
// 	conn, ok :=msg.Manager.ChatRoomMember[uid]

// 	conn, err := (&websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
// 			return true
// 		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

// 	// conn, err := websocket.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Logger.Error("websocket建立连接错误", err)
// 		http.NotFound(c.Writer, c.Request)
// 		return
// 	}
// 	// 创建一个用户实例
// 	client := &msg.Client{
// 		Uid:    uid,
// 		Socket: conn,
// 		Send:   make(chan []byte),
// 	}
// 	// fmt.Println("绑定client", client)
// 	// 用户注册到用户管理上
// 	msg.Manager.Register <- client
// 	go client.Read()
// 	go client.Write()

// }
