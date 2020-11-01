package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// EventData - A data structure that contains information about the current chat event.
type EventData struct {
	Event       string    `json:"event"`
	Body        string    `json:"body"`
	UserCount   int32     `json:"userCount"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
	Auth        string    `json:"auth"`
}

// Users - A map containing all the connected users.
var Users sync.Map

// UserCount - Total count of connected users.
var UserCount int32

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getUserName(connection *websocket.Conn) string {
	value, _ := Users.Load(connection)
	user := value.(*User)
	return user.Name
}

func removeUser(user *User) {
	Users.Delete(user)
	atomic.AddInt32(&UserCount, -1)
}

// SendToAll - sends the body string data to all connected clients
func SendToAll(body string, name string, eventType string) {
	log.Println("sendToAll(): " + body)
	response := EventData{Event: eventType, Body: body, UserCount: UserCount, Name: name, CreatedDate: time.Now()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Print("sendToAll():", err)
	}
	if eventType == EventMessage {
		UpdateChatHistory(jsonResponse)
	}
	Users.Range(func(key, value interface{}) bool {
		var userValue = value.(*User)
		if err := userValue.write(websocket.TextMessage, jsonResponse); err != nil {
			log.Print("sendToAll():", err)
		}
		return true
	})
}

// SendToOne - sends the body string data to a parameter defined client
func SendToOne(body string, user *User, eventType string) {
	log.Println("sendToOne(): " + body)
	response := EventData{Event: eventType, Body: body,
		UserCount: UserCount, Name: "", CreatedDate: time.Now()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Print("sendToOne():", err)
	}
	if err := user.write(websocket.TextMessage, jsonResponse); err != nil {
		log.Print("sendToOne():", err)
	}
}

// SendToOther - sends the body string data to all connected clients except the parameter given client
func SendToOther(body string, user *User, eventType string) {
	log.Print("sendToOther():", body)
	response := EventData{Event: eventType, Body: body, UserCount: UserCount, CreatedDate: time.Now()}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Print("sendToOther():", err)
	}
	if eventType == EventMessage {
		UpdateChatHistory(jsonResponse)
	}
	Users.Range(func(key, value interface{}) bool {
		userValue := value.(*User)
		if userValue != user {
			if err := userValue.write(websocket.TextMessage, jsonResponse); err != nil {
				log.Print("sendToOther():", err)
			}
		}
		return true
	})
}

func newChatConnection(connection *websocket.Conn, cookie string) {
	log.Print("chatRequest():", cookie)
	log.Print("chatRequest():", "Connection opened.")
	nano := strconv.Itoa(int(time.Now().UnixNano()))
	newUser := User{Name: "Anon" + nano, Connection: connection}
	Users.Store(&newUser, &newUser)
	atomic.AddInt32(&UserCount, 1)
	err := HandleJoin(&newUser)
	if err != nil {
		connection.Close()
		removeUser(&newUser)
		log.Print("newChatConnection():", err)
	} else {
		go reader(&newUser)
	}
}

func reader(user *User) {
	var readerError error
	defer func() {
		log.Print("reader():", readerError)
		user.Connection.Close()
		key, _ := Users.Load(user)
		user := key.(*User)
		removeUser(user)
		SendToAll(user.Name+" has left the chatroom.", "", EventNotification)
	}()
	for {
		var EventData EventData
		messageType, message, readerError := user.Connection.ReadMessage()
		if readerError != nil {
			return
		}
		if messageType == websocket.TextMessage {
			readerError = json.Unmarshal(message, &EventData)
			if len(EventData.Auth) > 0 {
				gatewayRes, err := refreshToken(EventData.Auth)
				if err != nil {
					SendToOne("Session error. Disconnected from chat.", user, EventNotification)
					readerError = err
				}
				EventData.Auth = gatewayRes.Token
			}
			if readerError != nil {
				return
			}
			switch EventData.Event {
			case EventTyping:
				readerError = HandleTypingEvent(EventData.Body, user)
			case EventMessage:
				readerError = HandleMessageEvent(EventData.Body, user)
			case EventLogin:
				readerError = HandleLoginEvent(EventData.Body, user)
			case EventNameChange:
				readerError = HandleNameChangeEvent(EventData.Body, user, EventData.Auth)
			}
			if readerError != nil {
				return
			}
		}
	}
}

// ChatRequest - A chat request.
func ChatRequest(responseWriter http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		allowedOrigin, found := os.LookupEnv("ALLOWED_ORIGIN")
		if found {
			return r.Header.Get("Origin") == "http://"+allowedOrigin ||
				r.Header.Get("Origin") == "https://"+allowedOrigin ||
				r.Header.Get("Origin") == "https://www."+allowedOrigin ||
				r.Header.Get("Origin") == "http://www."+allowedOrigin
		}
		return true
	}
	wsConnection, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Print("ChatRequest():", err)
	} else {
		newChatConnection(wsConnection, request.Header.Get("Cookie"))
	}
}
