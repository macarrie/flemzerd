package telegram

import (
	"fmt"
	"strconv"
	"time"

	notifier_helper "github.com/macarrie/flemzerd/helpers/notifiers"
	log "github.com/macarrie/flemzerd/logging"
	. "github.com/macarrie/flemzerd/objects"

	"github.com/macarrie/flemzerd/configuration"
	"github.com/macarrie/flemzerd/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type TelegramNotifier struct {
	Client   *tgbotapi.BotAPI
	AuthCode int
	ChatID   int64
}

var module Module

func New() (t *TelegramNotifier, err error) {
	t = &TelegramNotifier{}

	module = Module{
		Name: t.GetName(),
		Type: "notifier",
		Status: ModuleStatus{
			Alive:   true,
			Message: "",
		},
	}

	bot, err := tgbotapi.NewBotAPI(configuration.TELEGRAM_BOT_TOKEN)
	if err != nil {
		return nil, errors.Wrap(err, "error when creating Telegram bot")
	}

	bot.Debug = true
	t.Client = bot

	chat_id := db.Session.TelegramChatID
	if chat_id != 0 {
		t.ChatID = chat_id
	} else {
		log.Warning("No Telegram chat ID found. User will need to authorize access to Telegram in UI")
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
	if t.ChatID == 0 {
		module.Status.Alive = false
		module.Status.Message = "No chat ID found. Begin setup process in UI to enable Telegram notifications"
		return module, errors.New(module.Status.Message)
	}

	module.Status.Alive = true
	module.Status.Message = ""

	return module, nil
}

func (t *TelegramNotifier) GetName() string {
	return "telegram"
}

func (t *TelegramNotifier) Send(notif Notification) error {
	log.Debug("Sending Telegram notification")

	if configuration.TELEGRAM_BOT_TOKEN == "" {
		return errors.New("Telegram bot token not found")
	}
	if t.Client == nil {
		return errors.New("Could not contact Telegram bot to send notification")
	}

	title, content, err := notifier_helper.GetNotificationText(notif)
	if err != nil {
		return err
	}

	notif_content := fmt.Sprintf("%s: \n%s", title, content)
	msg := tgbotapi.NewMessage(t.ChatID, notif_content)
	if _, err = t.Client.Send(msg); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Unable to send notification")
		return errors.Wrap(err, "cannot send telegram notification")
	}

	return nil
}

func (t *TelegramNotifier) Auth() {
	if t.AuthCode != 0 {
		log.Debug("Telegram auth process already in progress. Skipping")
		return
	}
	if t.ChatID != 0 {
		log.Debug("Telegram setup already completed")
		return
	}

	errorsChannel := make(chan error, 1)
	doneChannel := make(chan bool, 1)

	pollTicker := time.NewTicker(3 * time.Second)
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, _ := t.Client.GetUpdatesChan(u)
		time.Sleep(time.Millisecond * 500)
		updates.Clear()

		for update := range updates {
			if update.Message == nil {
				continue
			}

			fmt.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if t.AuthCode != 0 && update.Message.Text == strconv.Itoa(t.AuthCode) {
				t.ChatID = update.Message.Chat.ID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Telegram linked sucessfully to flemzerd. Notifications will be sent through Telegram when the Telegram notifier is enabled in flemzerd.")
				msg.ReplyToMessageID = update.Message.MessageID
				if t.ChatID == 0 {
					t.ChatID = update.Message.Chat.ID
				}

				t.Client.Send(msg)
				doneChannel <- true
			}

		}
		<-pollTicker.C
	}()

	go func() {
		time.Sleep(120 * time.Second)
		errorsChannel <- errors.New("[Telegram auth] Auth code expired")
		t.AuthCode = 0
		doneChannel <- true
	}()

	<-doneChannel
	pollTicker.Stop()

	select {
	default:
		t.AuthCode = 0
		db.SaveTelegramChatID(t.ChatID)

		return
	}
}
