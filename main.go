package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"joonas.ninja-chat/routes/ws"
	"fmt"
)

func initEnvFile() {
	var err = godotenv.Load("app.env");
	if err != nil {
		panic("Error loading app.env file. Please create one next to me.");
	}
	fmt.Println("initEnvFile(): Loaded envs.");
}

func initRoutes() {
	http.HandleFunc("/api/v1/ws/chat", ws.ChatRequest);
}

func main() {
	initEnvFile();
	var err = http.ListenAndServe(":"+os.Getenv("PORT"), nil);
	if (err != nil) {
		panic(err);
	} else {
		initRoutes();
		fmt.Println("Päällä on.");
	}
}
