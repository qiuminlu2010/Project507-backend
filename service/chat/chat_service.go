package chat_service

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	FromUid  int    `json:"from_uid" form:"from_uid"`
	ToUid    int    `json:"to_uid" form:"to_uid"`
	Username string `json:"username" form:"username"`
	Avator   string `json:"avator" form:"avator"`
	Content  string `json:"content" form:"content"`
	Type     int    `json:"type"`
	Ctime    int    `json:"ctime"`
}

// 用户类
type Client struct {
	ID     int
	SendID int
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播类，包括广播内容和源用户
type Broadcast struct {
	Client *Client
	Msg    []byte
	Type   int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
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
			log.Println("数据格式不正确", err)
			c.close()
			break
		}

		fmt.Println(c.ID, "发送消息", msg.Content)
		//TODO: 保存数据库
		Manager.Broadcast <- &Broadcast{
			Client: c,
			Msg:    []byte(msg.Content),
		}

	}
}
