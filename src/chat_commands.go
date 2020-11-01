package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// HandleChannelCommand - Tyyppiä /channel <lisäkomento tähän>
func HandleChannelCommand(commands []string, user *User) {
	var subCommand = commands[2]
	switch subCommand {
	case SubCommandCreate:
		
	}
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
