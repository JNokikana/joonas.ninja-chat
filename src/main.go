package main

import (
	"net/http"
	"os"

	"fmt"

	"github.com/joho/godotenv"
)

func initEnvFile() {
	var err = godotenv.Load("app.env")
	if err != nil {
		panic("Error loading app.env file. Please create one next to me.")
	}
	fmt.Println("initEnvFile(): Loaded envs.")
}

func initRoutes() {
	http.HandleFunc("/api/v1/ws/chat", ChatRequest)
	fmt.Println("initRoutes(): Routes initialized.")
}

func main() {
	initEnvFile()
	initRoutes()
	fmt.Println("main(): Starting server...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}