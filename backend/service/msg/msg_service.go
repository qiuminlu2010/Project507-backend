package msg_service

import (
	"qiu/backend/model"
	"qiu/backend/pkg/e"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/redis"
	"qiu/backend/pkg/util"
	cache "qiu/backend/service/cache"
	param "qiu/backend/service/param"
	user "qiu/backend/service/user"
	"strconv"

	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	FromUid   int    `json:"from_uid" form:"from_uid"`
	ToUid     int    `json:"to_uid" form:"to_uid"`
	Username  string `json:"username" form:"username"`
	Avatar    string `json:"avatar" form:"avatar"`
	Content   string `json:"content" form:"content"`
	ImageUrl  string `json:"image_url" form:"image_url"`
	Type      int    `json:"type"`
	CreatedOn int64  `json:"created_on"`
}

type ReplyMessage struct {
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct {
	Uid int
	// ToUid   int
	Socket *websocket.Conn
	Send   chan []byte
}

type Broadcast struct {
	Client *Client
	Msg    *Message
	Type   int
}

type ClientManager struct {
	Clients        map[int]*Client
	Broadcast      chan *Broadcast
	Reply          chan *Client
	Register       chan *Client
	Unregister     chan *Client
	ChatRoomMember map[int]bool
}

var (
	Manager     *ClientManager
	ChatRoomMsg chan *Message
	SystemMsg   chan *Message
	// LikeArticleChannel chan *param.ArticleLikeParams
	// LikeCommentChannel chan *param.LikeCommentParams

)

func Setup() {
	Manager = &ClientManager{
		Clients:    make(map[int]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Reply:      make(chan *Client),
		Unregister: make(chan *Client),
	}
	ChatRoomMsg = make(chan *Message)
	SystemMsg = make(chan *Message)
	// LikeArticleChannel = make(chan *param.ArticleLikeParams)
	// LikeCommentChannel = make(chan *param.LikeCommentParams)
	// initChannelManager()
	go Manager.Listen()
}

func initChannelManager() {

	// pubLikeArticle := redis.Subscribe(e.CHANNEL_LIKEARTICLE)
	// defer pubLikeArticle.Close()
	// go func() {
	// 	for {
	// 		msg, err := pubLikeArticle.ReceiveMessage(ctx)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		log.Logger.Info(msg.Channel, msg.Payload)
	// 	}
	// }()
	// LikeArticleChannel = pubLikeArticle.Channel()

	// pubLikeComment := redis.Subscribe(e.CHANNEL_LIKECOMMENT)
	// defer pubLikeComment.Close()
	// LikeCommentChannel = pubLikeComment.Channel()
}
func (c *Client) close() {
	Manager.Unregister <- c
	_ = c.Socket.Close()
}

func (c *Client) Read() {
	defer c.close()
	for {
		c.Socket.PongHandler()
		msg := new(Message)
		err := c.Socket.ReadJSON(&msg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Logger.Error("数据格式不正确", err)
			c.close()
			break
		}

		userInfo := user.GetUserCache(c.Uid)
		msg.FromUid = c.Uid
		msg.CreatedOn = time.Now().Unix()
		msg.Username = userInfo.Name
		msg.Avatar = userInfo.Avatar

		msgModel := model.Message{
			FromUid:  c.Uid,
			ToUid:    msg.ToUid,
			Content:  msg.Content,
			ImageUrl: msg.ImageUrl,
		}
		log.Logger.Debug("[发送消息]", " 用户id: ", msg.FromUid, " 接收用户id ", msg.ToUid)
		if msg.ToUid > 0 {
			//私信
			Manager.Broadcast <- &Broadcast{
				Client: c,
				Msg:    msg,
			}
			//TODO: 保存数据库
			if err := model.SaveMessage(&msgModel); err != nil {
				panic(err)
			}

			//未读消息
			key := cache.GetModelFieldKey(e.CACHE_USER, uint(msg.ToUid), e.CACHE_UNREAD_MSG)
			redis.HashIncr(key, strconv.Itoa(msg.FromUid), 1)
			redis.Expire(key, e.DURATION_USER_SESSIONS)
		} else {
			//聊天室
			ChatRoomMsg <- msg
		}

	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
		// close(c.Send)
	}()
	for {
		msg, ok := <-c.Send
		if !ok {
			_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		_ = c.Socket.WriteMessage(websocket.TextMessage, msg)

	}
}
func (manager *ClientManager) Close() {
	log.Logger.Info("<---关闭管道通信--->")
	close(manager.Broadcast)
	close(manager.Register)
	close(manager.Unregister)
	close(manager.Reply)
}

func (manager *ClientManager) Listen() {
	// defer manager.Close()
	for {
		log.Logger.Info("<---监听WebSocket通信--->")
		select {
		case conn := <-manager.Register: // 建立连接
			// 关闭旧连接
			oldClient, ok := manager.Clients[conn.Uid]
			if ok {
				close(oldClient.Send)
				delete(manager.Clients, conn.Uid)
			}

			manager.Clients[conn.Uid] = conn

			replyMsg := &ReplyMessage{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			log.Logger.Infof("[Chat] 建立新连接: Uid%v", conn.Uid)
			replyMsgBytes, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
		case conn := <-manager.Unregister: // 断开连接
			log.Logger.Info("[Chat] 断开连接: Uid%v", conn.Uid)
			if _, ok := manager.Clients[conn.Uid]; ok {
				replyMsg := &ReplyMessage{
					Code:    e.WebsocketEnd,
					Content: "连接已断开",
				}
				replyMsgBytes, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
				// close(conn.Send)
				delete(manager.Clients, conn.Uid)
			}
		//广播信息
		case broadcast := <-manager.Broadcast:
			msg := broadcast.Msg
			ToUid := msg.ToUid
			sendClient, ok := manager.Clients[ToUid]
			msgBytes, _ := json.Marshal(msg)

			if ok {
				select {
				case sendClient.Send <- msgBytes:
					replyMsg := &ReplyMessage{
						Code:    e.WebsocketOnlineReply,
						Content: "对方在线应答",
					}
					replyMsgBytes, _ := json.Marshal(replyMsg)
					_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
				default:
					close(sendClient.Send)
					delete(Manager.Clients, sendClient.Uid)
				}
			} else {
				replyMsg := &ReplyMessage{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				replyMsgBytes, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
			}
			// 保存发送方最近会话
			model.SaveSession(&model.MessageSession{
				Uid:       msg.FromUid,
				SessionId: msg.ToUid,
			})
			key1 := cache.GetModelFieldKey(e.CACHE_USER, uint(msg.FromUid), e.CACHE_SESSIONS)
			if redis.Exists(key1) != 0 {
				redis.ZAdd(key1, float64(time.Now().Unix()), msg.ToUid)
			}
			// 保存接收方最近会话
			model.SaveSession(&model.MessageSession{
				Uid:       msg.ToUid,
				SessionId: msg.FromUid,
			})
			key2 := cache.GetModelFieldKey(e.CACHE_USER, uint(msg.ToUid), e.CACHE_SESSIONS)
			if redis.Exists(key2) != 0 {
				redis.ZAdd(key2, float64(time.Now().Unix()), msg.FromUid)
			}

		case roomMsg := <-ChatRoomMsg:
			msgBytes, _ := json.Marshal(roomMsg)
			log.Logger.Debug("[聊天室消息] ", "用户id: ", roomMsg.FromUid, " ", roomMsg.Content)
			for _, sendClient := range manager.Clients {
				select {
				case sendClient.Send <- msgBytes:
				default:
					close(sendClient.Send)
					delete(Manager.Clients, sendClient.Uid)
				}
			}
		case msg := <-SystemMsg:
			ToUid := msg.ToUid
			sendClient, ok := manager.Clients[ToUid]
			msgBytes, _ := json.Marshal(msg)
			if ok {
				select {
				case sendClient.Send <- msgBytes:
				default:
					close(sendClient.Send)
					delete(Manager.Clients, sendClient.Uid)
				}
			}
			log.Logger.Info("[推送消息]：", ToUid, msg.Content)

			msgModel := model.Message{
				FromUid:   msg.FromUid,
				ToUid:     msg.ToUid,
				Content:   msg.Content,
				CreatedOn: int(msg.CreatedOn),
			}
			//保存数据库
			if err := model.SaveMessage(&msgModel); err != nil {
				panic(err)
			}

			//记录未读消息
			key := cache.GetModelFieldKey(e.CACHE_USER, uint(msg.ToUid), e.CACHE_UNREAD_MSG)
			redis.HashIncr(key, strconv.Itoa(msg.FromUid), 1)
			redis.Expire(key, e.DURATION_USER_SESSIONS)

			// 保存接收方最近会话
			model.SaveSession(&model.MessageSession{
				Uid:       msg.ToUid,
				SessionId: msg.FromUid,
			})
			key2 := cache.GetModelFieldKey(e.CACHE_USER, uint(msg.ToUid), e.CACHE_SESSIONS)
			if redis.Exists(key2) != 0 {
				redis.ZAdd(key2, float64(time.Now().Unix()), msg.FromUid)
			}

		}
	}
}

func setSession(uid, pageNum, pageSize int) error {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_SESSIONS)
	// var sessions  []*model.MessageSession

	sessions, err := model.GetSession(uid, pageNum, pageSize)

	if err != nil {
		return nil
	}

	for _, session := range sessions {
		redis.ZAdd(key, float64(session.ModifiedOn), session.SessionId)
	}
	redis.Expire(key, e.DURATION_USER_SESSIONS)

	return nil

}

func setAllSessions(uid int) error {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_SESSIONS)
	// var sessions  []*model.MessageSession

	sessions, err := model.GetAllSessions(uid)

	if err != nil {
		return nil
	}

	for _, session := range sessions {
		redis.ZAdd(key, float64(session.ModifiedOn), session.SessionId)
	}
	redis.Expire(key, e.DURATION_USER_SESSIONS)

	return nil

}

func GetSession(uid, pageNum, pageSize int) ([]*model.SessionInfo, error) {

	key := cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_SESSIONS)

	// 分页缓存

	// if redis.Exists(key) == 0 {
	// 	if err := setSession(uid, 0, pageSize*(pageNum+1)); err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	cnt := redis.ZCard(key)
	// 	if (int64)(pageNum*pageSize) >= cnt {
	// 		offset := int(cnt)/pageSize + int(cnt)%(pageSize)
	// 		err := setSession(uid, offset, pageSize*(pageNum+1)-int(cnt))
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	// 一次性缓存所有
	if redis.Exists(key) == 0 {
		if err := setAllSessions(uid); err != nil {
			return nil, err
		}
	}
	value := redis.ZRevRange(key, int64(pageNum), int64(pageSize))
	userIds, err := util.StringsToInts(value)

	if err != nil {
		return nil, err
	}
	key = cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_UNREAD_MSG)
	var sessionInfos []*model.SessionInfo
	userInfos := user.GetUsersCache(userIds)
	for _, userInfo := range userInfos {
		sx := redis.HashGet(key, strconv.Itoa(int(userInfo.ID)))
		var unread int
		if len(sx) == 0 {
			unread = 0
		} else {
			unread, err = strconv.Atoi(sx)
			if err != nil {
				panic(err)
			}
		}
		messages, err := model.GetMessage(uid, int(userInfo.ID), pageNum, pageSize)
		if err != nil {
			panic(err)
		}
		sessionInfo := model.SessionInfo{
			UserBase: *userInfo,
			Unread:   unread,
			Messages: messages,
		}
		sessionInfos = append(sessionInfos, &sessionInfo)
	}
	return sessionInfos, nil
}

func CheckToken(token string, uid int) bool {
	if token == "" {
		return false
	}
	key := cache.GetModelFieldKey("user", uint(uid), "token")
	adminKey := cache.GetModelFieldKey("user", 2, "token")
	if redis.Exists(adminKey) != 0 {
		admin_token := redis.Get(adminKey)

		if admin_token == token {
			return true
		}
	}

	if redis.Exists(key) != 0 {
		cache_token := redis.Get(key)
		// fmt.Println("AdminToken", admin_token, token == admin_token)
		return token == cache_token
	}
	return false
}

func UpdateUnReadMessage(uid int, sessionId int) {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_UNREAD_MSG)
	redis.HashDel(key, strconv.Itoa(sessionId))

}
func GetMessages(params *param.MessageGetParams) ([]*model.Message, error) {
	// UpdateUnReadMessage(params.FromUid, params.ToUid)
	return model.GetMessage(params.FromUid, params.ToUid, params.Offset, params.Limit)
}
