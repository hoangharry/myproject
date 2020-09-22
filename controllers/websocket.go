package controllers

import(
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"encoding/json"
	"net/http"
    "github.com/astaxie/beego"
    "myproject/models"
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
    socket *websocket.Conn
    send   chan []byte
}

type Message struct {
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Content   string `json:"content,omitempty"`
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

func (manager *ClientManager) Send(message []byte, ignore *Client) {
    for conn := range manager.clients {
        if conn != ignore {
            conn.send <- message
        }
    }
}

// func (c *Client) Read() {
//     defer func() {
//         manager.unregister <- c
//         c.socket.Close()
//     }()

//     for {
//         _, message, err := c.socket.ReadMessage()
//         if err != nil {
//             manager.unregister <- c
//             c.socket.Close()
//             break
//         }
//         jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
//         manager.broadcast <- jsonMessage
//     }
// }

// func (this *WebsocketController) WsAPI(){
//     ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
//     if _, ok := err.(websocket.HandshakeError); ok {
//         http.Error(this.Ctx.ResponseWriter, "not a websocket handshake", 400)
//         return 
//     } else if err != nil {
//         beego.Error("cannot setup websocket connection", err)
//         return
//     }

// }

func BroadcastWs(img models.Image){
    img.Evidence = img.Evidence[2: len(img.Evidence) - 1]
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

// func WsPage(res http.ResponseWriter, req *http.Request) {
func (this* WebsocketController) WsPage() {
    res := this.Ctx.ResponseWriter
    req := this.Ctx.Request
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if error != nil {
        http.NotFound(res, req)
        return
    }
    client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}

    manager.register <- client

    // go client.Read()
    go client.Write()
}

func init() {
    fmt.Println("Starting application...")
    go manager.Start()
    // http.HandleFunc("/ws", WsPage)
    // http.ListenAndServe(":8000", nil)
}

// WebSocketController handles WebSocket requests.
// type WebSocketController struct {
// 	beego.Controller
// }

// // Get method handles GET requests for WebSocketController.
// func (this *WebSocketController) Get() {
// 	// Safe check.
// 	uname := this.GetString("uname")
// 	if len(uname) == 0 {
// 		this.Redirect("/", 302)
// 		return
// 	}

// 	this.TplName = "websocket.html"
// 	this.Data["IsWebSocket"] = true
// 	this.Data["UserName"] = uname
// }

// Join method handles WebSocket requests for WebSocketController.
// func (this *WebSocketController) Join() {
	// uname := this.GetString("uname")
	// if len(uname) == 0 {
	// 	this.Redirect("/", 302)
	// 	return
	// }

	// Upgrade from http request to WebSocket.
// 	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
// 	if _, ok := err.(websocket.HandshakeError); ok {
// 		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
// 		return
// 	} else if err != nil {
// 		beego.Error("Cannot setup WebSocket connection:", err)
// 		return
// 	}

// 	// Join chat room.
// 	Join(uname, ws)
// 	defer Leave(uname)

// 	// Message receive loop.
// 	for {
// 		_, p, err := ws.ReadMessage()
// 		if err != nil {
// 			return
// 		}
// 		publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
// 	}
// }

// // broadcastWebSocket broadcasts messages to WebSocket users.
// func broadcastWebSocket(event models.Event) {
// 	data, err := json.Marshal(event)
// 	if err != nil {
// 		beego.Error("Fail to marshal event:", err)
// 		return
// 	}

// 	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
// 		// Immediately send event to WebSocket users.
// 		ws := sub.Value.(Subscriber).Conn
// 		if ws != nil {
// 			if ws.WriteMessage(websocket.TextMessage, data) != nil {
// 				// User disconnected.
// 				unsubscribe <- sub.Value.(Subscriber).Name
// 			}
// 		}
// 	}
// }
