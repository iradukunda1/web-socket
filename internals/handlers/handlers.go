package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var WsChann = make(chan WsPayload)
var clients = make(map[WebScoketConnection]string)

func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}

}

type WebScoketConnection struct {
	*websocket.Conn
}

//WsHandler is the handler that upgrade the connection to websocket
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("client connected to websocket")

	var res WsJsonResponse

	res.Message = `<em> <small>Hello, welcome to the websocket</small></em>`

	conn := WebScoketConnection{Conn: ws}

	clients[conn] = ""

	err = ws.WriteJSON(res)
	if err != nil {
		log.Println(err)
	}

	go ListenFormWs(&conn)

}

func ListenFormWs(conn *WebScoketConnection) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("err", fmt.Sprintf("%v", err))
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothings
		} else {
			payload.Conn = *conn
			WsChann <- payload
		}
	}
}

func ListenToWsChannel() {
	var res WsJsonResponse
	for {
		evt := <-WsChann

		switch evt.Action {
		case "username":
			//get a list of all clients and send it back via broadcast
			clients[evt.Conn] = evt.UserName
			users := getClientList()
			res.Action = "list_users"
			res.ConnectedUser = users
			boardCastToAll(res)
		case "left":
			res.Action = "list_users"
			delete(clients, evt.Conn)
			users := getClientList()
			res.ConnectedUser = users
			boardCastToAll(res)
		case "boardcast":
			res.Action = "boardcast"
			res.Message = fmt.Sprintf("<strong>%s</strong>: %s", evt.UserName, evt.Message)
			boardCastToAll(res)
		}
	}
}

func getClientList() []string {
	var userList []string
	for _, client := range clients {
		if client != "" {
			userList = append(userList, client)
		}
	}

	sort.Strings(userList)

	return userList
}

func boardCastToAll(res WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(res)
		if err != nil {
			log.Println("web socket close")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

func renderPage(w http.ResponseWriter, tmp string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmp)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
