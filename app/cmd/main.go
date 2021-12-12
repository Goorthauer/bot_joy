package main

import (
	"bot_joy/app/internal/model"
	"log"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().Unix())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Бот включен")
	redisClient, err := model.NewRedis()
	model.TelegramBot(redisClient)
}
