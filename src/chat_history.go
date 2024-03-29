package main

import (
	"encoding/json"
	"log"
)

type ChatHistory struct {
	Body      []EventData `json:"history"`
	UserCount int32       `json:"userCount"`
	Event     string      `json:"event"`
}
 
// updateChatHistory - Adds the parameter defined chat history entry to chat history
func updateChatHistory(jsonResponse []byte) {
	go apiRequest("POST", newApiRequestOptions(&apiRequestOptions{payload: jsonResponse}), "CHAT_HISTORY_URL", nil)
}

func getChatHistory(channelId string) ChatHistory {
	res, err := apiRequest("GET", newApiRequestOptions(&apiRequestOptions{queryString: "?channelId=" + channelId}), "CHAT_HISTORY_URL", nil)
	if err != nil {
		log.Print("getChatHistory():", err)
		return ChatHistory{}
	}
	var eventData []EventData
	if err := json.Unmarshal(res, &eventData); err != nil {
		log.Print("getChatHistory():", err)
		return ChatHistory{}
	}
	return ChatHistory{Event: EventChatHistory, Body: eventData, UserCount: UserCount}
}
