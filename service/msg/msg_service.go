package msg_service

import (
	"encoding/json"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	log "qiu/blog/pkg/logging"
	user "qiu/blog/service/user"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	FromUid  int    `json:"from_uid" form:"from_uid"`
	ToUid    int    `json:"to_uid" form:"to_uid"`
	Username string `json:"username" form:"username"`
	Avator   string `json:"avator" form:"avator"`
	Content  string `json:"content" form:"content"`
	ImageUrl string `json:"image_url" form:"image_url"`
	Type     int    `json:"type"`
	Ctime    int64  `json:"ctime"`
}

type ReplyMessage struct {
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户类
type Client struct {
	FromUid int
	ToUid   int
	Socket  *websocket.Conn
	Send    chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client *Client
	// Msg    []byte
	Msg  *Message
	Type int
}

// 用户管理
type ClientManager struct {
	Clients    map[int]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager *ClientManager

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
			log.Error("数据格式不正确", err)
			c.close()
			break
		}
		// if redis.Exists()
		userInfo := user.GetUserCache(c.FromUid)
		msg.FromUid = c.FromUid
		msg.ToUid = c.ToUid
		msg.Ctime = time.Now().Unix()
		msg.Username = userInfo.Name
		msg.Avator = userInfo.Avator
		log.Info(c.FromUid, "发送消息", msg)
		//TODO: 保存数据库
		msgModel := model.Message{
			FromUid:  c.FromUid,
			ToUid:    c.ToUid,
			Content:  msg.Content,
			ImageUrl: msg.ImageUrl,
		}
		if err := model.SaveMessage(&msgModel); err != nil {
			panic(err)
		}
		Manager.Broadcast <- &Broadcast{
			Client: c,
			Msg:    msg,
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
		// log.Println(c.FromUid, "接受消息:", msg)
		_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		// replyMsg := &ReplyMessage{Message: *msg}
		// log.Println(c.FromUid, "接受消息:", replyMsg)
		// msgBytes, _ := json.Marshal(replyMsg)
		// _ = c.Socket.WriteMessage(websocket.TextMessage, msgBytes)

	}
}
func (manager *ClientManager) Close() {
	log.Info("<---关闭管道通信--->")
	close(manager.Broadcast)
	close(manager.Register)
	close(manager.Unregister)
	close(manager.Reply)
}

func (manager *ClientManager) Listen() {
	// defer manager.Close()
	for {
		log.Info("<---监听WebSocket通信--->")
		select {
		case conn := <-manager.Register: // 建立连接
			manager.Clients[conn.FromUid] = conn
			replyMsg := &ReplyMessage{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			log.Info("[Chat] 建立新连接: Uid%v", conn.FromUid)
			replyMsgBytes, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
		case conn := <-manager.Unregister: // 断开连接
			log.Info("[Chat] 断开连接: Uid%v", conn.FromUid)
			if _, ok := manager.Clients[conn.FromUid]; ok {
				replyMsg := &ReplyMessage{
					Code:    e.WebsocketEnd,
					Content: "连接已断开",
				}
				replyMsgBytes, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
				close(conn.Send)
				delete(manager.Clients, conn.FromUid)
			}
		//广播信息
		case broadcast := <-manager.Broadcast:
			msg := broadcast.Msg
			ToUid := broadcast.Client.ToUid
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
					delete(Manager.Clients, sendClient.FromUid)
				}
			} else {
				replyMsg := &ReplyMessage{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				replyMsgBytes, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, replyMsgBytes)
				// err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				// if err != nil {
				// 	fmt.Println("InsertOneMsg Err", err)
				// }
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
	go Manager.Listen()
}
