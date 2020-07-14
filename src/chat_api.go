package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type chatHistory struct {
	Body      []EventData `json:"history"`
	UserCount int32       `json:"userCount"`
	Event     string      `json:"event"`
}

type channelResponse struct {
	Body      []EventData `json:"history"`
	UserCount int32       `json:"userCount"`
	Event     string      `json:"event"`
}

// UpdateChatHistory - Adds the parameter defined chat history entry to chat history
func UpdateChatHistory(jsonResponse []byte) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("CHAT_HISTORY_URL"), bytes.NewBuffer(jsonResponse))
	if err != nil {
		log.Print("updateChatHistory():", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", `Basic `+
		base64.StdEncoding.EncodeToString([]byte(os.Getenv("APP_ID")+":"+os.Getenv("API_KEY"))))
	historyResponse, err := client.Do(req)
	if historyResponse != nil && historyResponse.Status != "200 OK" {
		log.Print("updateChatHistory():", "Error response "+historyResponse.Status)
	}
	if err != nil {
		log.Print("updateChatHistory():", err)
	}
	defer historyResponse.Body.Close()
}

// GetChatHistory - Returns the entire chat history.
func GetChatHistory() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", os.Getenv("CHAT_HISTORY_URL"), nil)
	req.Header.Add("Authorization", `Basic `+
		base64.StdEncoding.EncodeToString([]byte(os.Getenv("APP_ID")+":"+os.Getenv("API_KEY"))))
	historyResponse, err := client.Do(req)
	if historyResponse != nil && historyResponse.Status != "200 OK" {
		log.Print("getChatHistory():", "Error response "+historyResponse.Status)
		return nil
	}
	if err != nil {
		log.Print("getChatHistory():", err)
		return nil
	}
	defer historyResponse.Body.Close()
	body, err := ioutil.ReadAll(historyResponse.Body)
	if err != nil {
		log.Print("getChatHistory():", err)
		return nil
	}
	var historyArray []EventData
	err = json.Unmarshal(body, &historyArray)
	if err != nil {
		log.Print("getChatHistory():", err)
		return nil
	}
	response := chatHistory{Event: EventChatHistory, Body: historyArray, UserCount: UserCount}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Print("getChatHistory():", err)
		return nil
	}
	return jsonResponse
}

// CreateChatChannel - Creates a chat channel and adds the user with the corresponding email to the new channel. Requires user authentication.
func CreateChatChannel(channelName string, userEmail string) []byte {
	return nil
}