package model

import (
	joyConfig "bot_joy/config"
	"context"
	"fmt"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	queue *MessageQueue
)

func getCommand(conf *joyConfig.Config) string {
	response := "комадны\n"
	for i, res := range conf.Query {
		response += fmt.Sprintf("-----------\nБлок №%v\nЗапрос: %v\nТеги для поиска: %v\n-----------", i+1, res.Call, res.Response)
	}
	return fmt.Sprint(response)

}

func TelegramBot(r *RedisClient) {
	conf := joyConfig.New()
	b, err := tb.NewBot(tb.Settings{
		Token:  conf.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	queue = NewMessageQueue(b, 100*time.Millisecond)

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		ctx, cancel := context.WithTimeout(context.Background(), 33*time.Second)
		defer cancel()
		upperText := strings.ToUpper(m.Text)
		if strings.Contains(upperText, "/HELP") {
			queue.bot.Send(m.Chat, getCommand(conf))
		} else if strings.Contains(upperText, "/STATISTIC") {
			out, _ := r.getAll(ctx, strconv.Itoa(m.Sender.ID))
			queue.bot.Send(m.Chat, out)
		}
		tag := commandExist(conf.Query, upperText)
		if tag != "" {
			queue.MessageQueue <- Message{
				Chat:    m.Chat,
				Tag:     tag,
				Options: m,
			}
			log.Printf("Присвоен тег:%v. Специально для юзера %v", tag, m.Sender.Username)
			r.addMemory(ctx, tag, strconv.Itoa(m.Sender.ID))
		}
	})
	b.Start()
}

func cropImage(filename string) error {
	src, err := imaging.Open(filename)
	if err != nil {
		return err
	}
	var imageWith = src.Bounds().Dx()
	var imageHeight = src.Bounds().Dy()
	src = imaging.CropAnchor(src, imageWith, imageHeight-15, imaging.Top)
	image, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	jpeg.Encode(image, src, nil)
	return nil
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
