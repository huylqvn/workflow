package bot

import (
	"sync"
	"time"

	"github.com/huylqvn/httpserver/log"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// 3PL service: telegram
func NewBot(apiKey string) *Bot {
	log := log.Get()
	pref := tele.Settings{
		Token:  apiKey,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Error("tele.NewBot:", err)
		return nil
	}

	return &Bot{b, log}
}

type Bot struct {
	b   *tele.Bot
	log *logrus.Logger
}

var instance *Bot
var once sync.Once

func Get() *Bot {
	once.Do(func() {
		instance = NewBot("6189055077:AAFP714IMfX7SkOW95Jyzwn3NF4Fcf46m34")
	})
	return instance
}

func (b *Bot) SendGroup(id int64, msg string) error {
	options := &tele.SendOptions{
		ParseMode: tele.ModeHTML,
	}
	_, err := b.b.Send(&tele.Chat{
		ID: -id,
	}, msg, options)
	if err != nil {
		b.log.Error("SendError:", err)
	}
	// delay 100ms for next request, required telegram integration
	time.Sleep(100 * time.Millisecond)
	return err
}

func (b *Bot) GetBot() *tele.Bot {
	return b.b
}
