package ws;

import (
    "net/http"
    "fmt"
    "strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func reader(connection *websocket.Conn) {
    for {
        messageType, message, err := connection.ReadMessage();
        if (err != nil) {
            fmt.Println(err);
            return
        }
        fmt.Println(string(message));

        if err := connection.WriteMessage(messageType, message); err != nil {
            fmt.Println(err);
            return
        }
    }
}

func chatRequest(connection *websocket.Conn) {
    fmt.Println("chatRequest(): Connection opened.");
    reader(connection);
}

// WebsocketRequest - Handles websocket requests and conveys them to the handler depending on request path.
func WebsocketRequest(responseWriter http.ResponseWriter, request *http.Request) {
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    wsConnection, err := upgrader.Upgrade(responseWriter, request, nil);
    if (err != nil) {
        fmt.Println(err);
	}
    pathArray := strings.Split(request.RequestURI, "/");

    if (pathArray[len(pathArray) - 1] == "chat") {
        chatRequest(wsConnection);
    }
}