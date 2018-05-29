package telegram

import (
	"errors"
	"fmt"
	"time"

	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/macarrie/flemzerd/configuration"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramNotifier struct {
	Client *tgbotapi.BotAPI
	Token  string
	ChatID int64
}

var module Module

func New() (t *TelegramNotifier, err error) {
	t = &TelegramNotifier{}

	bot, err := tgbotapi.NewBotAPI(configuration.TELEGRAM_BOT_TOKEN)
	if err != nil {
		return nil, err
	}

	bot.Debug = true
	t.Client = bot

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, _ := bot.GetUpdatesChan(u)
		time.Sleep(time.Millisecond * 500)
		updates.Clear()

		for update := range updates {
			if update.Message == nil {
				continue
			}

			fmt.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			if t.ChatID == 0 {
				t.ChatID = update.Message.Chat.ID
			}

			bot.Send(msg)
		}
	}()

	module = Module{
		Name: "telegram",
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	return t, nil
}

func (t *TelegramNotifier) Status() (Module, error) {
	log.Debug("Checking telegram notifier status")

	if configuration.TELEGRAM_BOT_TOKEN == "" {
		module.Status.Alive = false
		module.Status.Message = "Telegram bot token not found"
		return module, errors.New(module.Status.Message)
	}
	if t.Client == nil {
		module.Status.Alive = false
		module.Status.Message = "Could not connect to telegram bot: no client"
		return module, errors.New(module.Status.Message)
	}

	// TODO

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (t *TelegramNotifier) GetName() string {
	return "telegram"
}

func (t *TelegramNotifier) Send(title, content string) error {
	log.WithFields(log.Fields{
		"title":   title,
		"content": content,
	}).Debug("Sending Telegram notification")

	if configuration.TELEGRAM_BOT_TOKEN == "" {
		return errors.New("Telegram bot token not found")
	}
	if t.Client == nil {
		return errors.New("Could not contact Telegram bot to send notification")
	}

	msg := tgbotapi.NewMessage(t.ChatID, content)
	_, err := t.Client.Send(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Unable to send notification")
		return err
	}

	return nil
}
