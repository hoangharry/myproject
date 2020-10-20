package controllers

import (
	"encoding/json"
	"fmt"
	"myproject/models"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type WebsocketController struct {
	beego.Controller
}
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	usr    string
	camID  string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

type UsrMes struct {
	Usr   string `json:"usr"`
	CamID string `json:"cam_id"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			fmt.Println("new hoomie")
			// jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected."})
			// manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				fmt.Println("homie out TT")
				// jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected."})
				// manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				var usr UsrMes
				if err := json.Unmarshal(message, &usr); err != nil {
					fmt.Println(err)
				}
				if conn.usr == usr.Usr && conn.camID == usr.CamID {
					select {
					case conn.send <- message:

					default:
						close(conn.send)
						delete(manager.clients, conn)
						fmt.Println("homie out TT")
					}
				}
			}
		}
	}
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func BroadcastWs(img models.Image) {
	img.Evidence = img.Evidence[2 : len(img.Evidence)-1]
	data, err := json.Marshal(img)
	if err != nil {
		beego.Error("fail", err)
		return
	}
	manager.broadcast <- data
}

func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (this *WebsocketController) WsPage() {
	res := this.Ctx.ResponseWriter
	req := this.Ctx.Request
	usr := this.Ctx.Input.Param(":usr")
	camID := this.Ctx.Input.Param(":camid")
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte), usr: usr, camID: camID}
	manager.register <- client
	go client.Write()
}

func init() {
	go manager.Start()
}
