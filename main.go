package main

import (
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Бот включен")
	telegramBot()
}
