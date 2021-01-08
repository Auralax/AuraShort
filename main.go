package main

import (
	"aurashort/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	loadEnv()
	server.StartServer()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
