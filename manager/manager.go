package manager

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	joyConfig "bot_joy/config"
	"bot_joy/internal/model"
)

type Manager struct {
	Redis *model.RedisClient
	Trace *model.TracingService
	Queue *model.MessageQueue
}

func New() *Manager {
	redisClient, err := model.NewRedis()
	if err != nil {
		log.Fatal("redis not connected %w", err)
	}
	return &Manager{Redis: redisClient, Trace: model.NewTrace()}
}

func (m *Manager) JoinBot() {
	conf := joyConfig.New()
	b, err := tb.NewBot(tb.Settings{
		Token:  conf.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	m.Queue = model.NewMessageQueue(b, 100*time.Millisecond)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(msg *tb.Message) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		upperText := strings.ToUpper(msg.Text)
		switch {
		case strings.Contains(upperText, "/HELP"):
			m.Queue.Bot.Send(msg.Chat, getCommand(conf), &tb.SendOptions{ReplyTo: msg})
		case strings.Contains(upperText, "/STATISTIC"):
			out, _ := m.Redis.GetAll(ctx, strconv.Itoa(msg.Sender.ID))
			m.Queue.Bot.Send(msg.Chat, out, &tb.SendOptions{ReplyTo: msg})
		}
		if tag := commandExist(conf.Query, upperText); tag != "" {
			spanTrace := model.NewSpan(ctx, msg.Text)
			m.Queue.MessageQueue <- model.Message{
				Chat:    msg.Chat,
				Tag:     tag,
				Options: msg,
				Trace:   spanTrace,
			}
			log.Printf("Присвоен тег:%v. Специально для юзера %v", tag, msg.Sender.Username)
			m.Redis.AddMemory(ctx, tag, strconv.Itoa(msg.Sender.ID))
		}
	})
	b.Start()
}

func commandExist(configs []joyConfig.QueryConfig, query string) string {
	for _, config := range configs {
		for _, item := range config.Call {
			if strings.Contains(query, strings.TrimSpace(item)) {
				return getTag(config.Response)
			}
		}
	}
	return ""
}

func getTag(reasons []string) string {
	return reasons[rand.Intn(len(reasons))]
}

func getCommand(conf *joyConfig.Config) string {
	response := "Команды:\n"
	for i, res := range conf.Query {
		response += fmt.Sprintf("-----------\nБлок №%v\nЗапрос: %v\nТеги для поиска: %v\n-----------", i+1, res.Call, res.Response)
	}
	return fmt.Sprint(response)

}
