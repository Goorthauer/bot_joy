package main

import (
	"log"
	"math/rand"
	"time"

	"bot_joy/manager"

	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().Unix())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Бот включен")
	manager.New().JoinBot()
}
