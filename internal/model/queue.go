package model

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
	Trace   *SpanLocal
}

type MessageQueue struct {
	Bot          *tb.Bot
	Timeout      time.Duration
	MessageQueue chan Message
}

func NewMessageQueue(bot *tb.Bot, timout time.Duration) *MessageQueue {
	queue := &MessageQueue{
		Bot:          bot,
		Timeout:      timout,
		MessageQueue: make(chan Message, 100),
	}
	go queue.QueueWorker()
	return queue
}

func (m *MessageQueue) QueueWorker() {
	for msg := range m.MessageQueue {
		m.send(msg)
	}
}

func (m *MessageQueue) send(msg Message) {
	var joyUrl string
	time.Sleep(m.Timeout)
	joyUrl, err := GetRandomBoobs(msg.Tag)
	msg.Trace.
		addInt("chat_id", msg.Options.Chat.ID).
		addString("chat_title", msg.Options.Chat.Title).
		addString("chat_username", msg.Options.Chat.Username).
		addString("joy_url", joyUrl).
		End()
	if err != nil {
		log.Println(err)
		_, err = m.Bot.Send(msg.Chat, "Ошибка получения изображения", &tb.SendOptions{ReplyTo: msg.Options.ReplyTo})
		if err != nil {
			log.Println(err)
		}
		return
	}

	filename, err := DownloadFile(joyUrl)
	if err != nil {
		log.Println(err)
		if _, err = m.Bot.Send(msg.Chat, "Ошибка распознавания изображения", &tb.SendOptions{ReplyTo: msg.Options.ReplyTo}); err != nil {
			log.Println(err)
		}
		return
	}
	if err = cropImage(filename); err != nil {
		log.Println(err)
		return
	}
	image := &tb.Photo{File: tb.FromDisk(filename)}
	if _, err = m.Bot.Send(msg.Chat, image, &tb.SendOptions{ReplyTo: msg.Options}); err != nil {
		log.Println(err)
		return
	}
	if err = os.Remove(filename); err != nil {
		log.Println(err)
		return
	}
}
