package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// HandleChannelCommand - For channel specific operations like list and create.
func HandleChannelCommand(commands []string, user *User) {
	var subCommand string
	var parameter string
	if len(commands) >= 2 {
		subCommand = commands[1]
		if subCommand != "" {
			switch parameter {
			case CommandChannelCreate:
				if len(commands) >= 3 {
					parameter = commands[2]
					CreateChatChannel(parameter, user.Email)
				}
			}
		}
	}
	log.Println("HandleChannelCommand(): ", "insufficient parameters for channel command.")
	return
}

// HandleUserCommand - sdgsdfg
func HandleUserCommand(commands []string, connection *websocket.Conn) {

}

// HandleWhoCommand - who is present in the current channel
func HandleWhoCommand(user *User) {
	var whoIsHere []string
	Users.Range(func(key, value interface{}) bool {
		v := value.(*User)
		whoIsHere = append(whoIsHere, v.Name)
		return true
	})
	jsonResponse, err := json.Marshal(whoIsHere)
	if err != nil {
		log.Printf("HandleWhoCommand(): ")
		log.Println(err)
		return
	}
	SendToOne(string(jsonResponse), user, EventWho)
}
