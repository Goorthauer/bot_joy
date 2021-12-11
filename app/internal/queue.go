package internal

import (
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Message struct {
	Chat    *tb.Chat
	Tag     string
	Options *tb.Message
}

type MessageQueue struct {
	bot          *tb.Bot
	timeout      time.Duration
	MessageQueue chan Message
}

func NewMessageQueue(bot *tb.Bot, timout time.Duration) *MessageQueue {
	queue := &MessageQueue{
		bot:          bot,
		timeout:      timout,
		MessageQueue: make(chan Message, 100),
	}
	go queue.QueueWorker()
	return queue
}

func (m *MessageQueue) QueueWorker() {
	for msg := range m.MessageQueue {
		time.Sleep(m.timeout)
		joyUrl, err := GetRandomBoobs(msg.Tag)
		if err != nil {
			log.Println(err)
			_, err = m.bot.Send(msg.Chat, "Ошибка получения изображения", &tb.SendOptions{ReplyTo: msg.Options.ReplyTo})
			if err != nil {
				log.Println(err)
			}
			continue
		}
		filename, err := DownloadFile(joyUrl)
		if err != nil {
			log.Println(err)
			_, err = m.bot.Send(msg.Chat, "Ошибка распознавания изображения", &tb.SendOptions{ReplyTo: msg.Options.ReplyTo})
			if err != nil {
				log.Println(err)
			}
			continue
		}
		err = cropImage(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		image := &tb.Photo{File: tb.FromDisk(filename)}
		_, err = m.bot.Send(msg.Chat, image, &tb.SendOptions{ReplyTo: msg.Options})
		if err != nil {
			log.Println(err)
			continue
		}
		err = os.Remove(filename)
		if err != nil {
			log.Println(err)
			continue
		}
	}

}
