package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func SessionRequest(responseWriter http.ResponseWriter, request *http.Request) {
	
}

func LoginRequest(responseWriter http.ResponseWriter, request *http.Request) {
	var loginRes loginDTO

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
	if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		err = json.Unmarshal(body, &loginRes)
		if err != nil {
			log.Print("getChatHistory():", err)
		}
		token, err := HandleLoginRequest(loginRes.Username, loginRes.Password)
		if (err != nil) {
			http.Error(responseWriter, "Unauthorized", http.StatusUnauthorized)
			return
		}
		isSecure, found := os.LookupEnv("IS_PROD")
		if (found && isSecure == "true"){
			isSecure = " secure;"
		} else {
			isSecure = ""
		}
		domain, found := os.LookupEnv("DOMAIN")
		if (!found){
			domain = ""
		}
		responseWriter.Header().Add("Set-Cookie", "session=" + token + "; httpOnly; sameSite=Strict; path=/;" + isSecure + "domain=" + domain + ";")
	} else {
		http.NotFound(responseWriter, request)
	}
}

// HandleLoginEvent - Handles the logic with user login.
func HandleLoginRequest(email string, password string) (loginToken string, err error) {
	var responseToken string

	if len(email) > 1 && len(password) > 1 {
		loginRes, loginError := loginRequest(email, password)
		if loginError != nil {
			return "", loginError
		} else {
			responseToken = loginRes.Token
		}
		log.Print("HandleLoginEvent():", "Login successful")
		return responseToken, nil
	}
	return "", nil
}