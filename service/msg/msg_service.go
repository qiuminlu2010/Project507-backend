package msg_service

import (
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	log "qiu/blog/pkg/logging"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"
	cache "qiu/blog/service/cache"
	user "qiu/blog/service/user"

	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	FromUid  int    `json:"from_uid" form:"from_uid"`
	ToUid    int    `json:"to_uid" form:"to_uid"`
	Username string `json:"username" form:"username"`
	Avatar   string `json:"avatar" form:"avatar"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image_url" form:"image_url"`
	Type     int    `json:"type"`
	Ctime    int64  `json:"ctime"`
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
	// Msg    []byte
	Msg  *Message
	Type int
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
)

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
		msg.Ctime = time.Now().Unix()
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
				close(conn.Send)
				delete(manager.Clients, conn.Uid)
			}
		//广播信息
		case broadcast := <-manager.Broadcast:
			msg := broadcast.Msg
			ToUid := msg.ToUid
			// flag := false // 默认对方不在线
			sendClient, ok := manager.Clients[ToUid]
			msgBytes, _ := json.Marshal(msg)
			// fromUid := broadcast.Client.FromUid
			if ok {
				select {
				case sendClient.Send <- msgBytes:
					replyMsg := &ReplyMessage{
						Code:    e.WebsocketOnlineReply,
						Content: "对方在线应答",
					}
					replyMsgBytes, _ := json.Marshal(replyMsg)
					_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
					// flag = true
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
			// 保存会话
			model.SaveSession(&model.MessageSession{
				Uid:       msg.FromUid,
				SessionId: msg.ToUid,
			})
			model.SaveSession(&model.MessageSession{
				Uid:       msg.ToUid,
				SessionId: msg.FromUid,
			})
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
		}
	}
}

func Setup() {
	Manager = &ClientManager{
		Clients:    make(map[int]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Reply:      make(chan *Client),
		Unregister: make(chan *Client),
	}
	ChatRoomMsg = make(chan *Message)
	go Manager.Listen()
}

func GetSession(uid, pageNum, pageSize int) ([]*model.UserBase, error) {
	key := cache.GetModelFieldKey(e.CACHE_USER, uint(uid), e.CACHE_SESSIONS)
	// var sessions  []*model.MessageSession
	if redis.Exists(key) == 0 {
		sessions, err := model.GetSession(uid, pageNum, pageSize)

		if err != nil {
			return nil, err
		}
		for _, session := range sessions {
			redis.ZAdd(key, float64(session.ModifiedOn), session.SessionId)
		}
		redis.Expire(key, e.DURATION_USER_SESSIONS)
		// return sessions, nil
	}
	value := redis.ZRevRange(key, int64(pageNum), int64(pageSize))
	userIds, err := util.StringsToInts(value)

	if err != nil {
		return nil, err
	}

	// followUsers := user.GetUsersCache(userIds)
	// json.Unmarshal(redis.ZRevRange(key, int64(pageNum), int64(pageSize)),&sessions)

	return user.GetUsersCache(userIds), nil
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
