package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WebSocketConnection struct {
	*websocket.Conn
}

type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Client connected to endpoint failed => ", err)
	}
	log.Println("Client is connected to endpoint..")

	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	// creating websocket connection for every connection
	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)

	if err != nil {
		log.Println("Error while connecting =====> ", err)
	}

	go ListenForWs(&conn)
}

func ListenToWsChannel() {
	var response WsJsonResponse

	for {
		evt := <-wsChan
		response.Action = "Got here"
		response.Message = fmt.Sprintf("Some chit chat %s", evt.Action)
		broadcastToAll(response)
	}
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error in websocket => ", err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error ", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//	Do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println("Error in render page ====>", err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println("Error in view execute ====>", err)
		return err
	}

	return nil
}
