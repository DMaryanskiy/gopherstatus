package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	StopReceivingUpdates()
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}

type DBInterface interface {
	GetUserByTelegram(username string) (storage.User, error)
	LatestResultsByID(limit int, userID uint) ([]storage.CheckResult, error)
}

type Bot struct {
	api TelegramAPI
	db  DBInterface
}

func NewBot(token string, db *storage.DB) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("failed to init bot api: %v", err)
		return nil, err
	}
	return &Bot{
		api: api,
		db:  db,
	}, nil
}

func (b *Bot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot listener...")
			b.api.StopReceivingUpdates()
			return
		case update, ok := <-updates:
			if !ok {
				log.Println("Updates channel closed.")
				return
			}
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}
			go b.handleCommand(update.Message)
		}
	}
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "status":
		b.handleStatus(msg)
	case "help":
		b.send(msg.Chat.ID, "/status — latest service check\n/help — this help message")
	default:
		b.send(msg.Chat.ID, "Unknown command. Try /help")
	}
}

func (b *Bot) handleStatus(msg *tgbotapi.Message) {
	user, err := b.db.GetUserByTelegram(msg.From.UserName)
	if err != nil {
		b.send(msg.Chat.ID, "Failed to fetch user from DB")
		return
	}
	results, err := b.db.LatestResultsByID(5, user.ID)
	if err != nil {
		b.send(msg.Chat.ID, "Failed to fetch results from DB")
		return
	}

	if len(results) == 0 {
		b.send(msg.Chat.ID, "No data available yet")
		return
	}

	text := "*Latest Service Checks:*\n\n"
	for _, r := range results {
		status := "❌ DOWN"
		if r.Online {
			status = "✅ UP"
		}
		text += fmt.Sprintf("*%s*: %s (%dms)\n", r.ServiceName, status, r.ResponseMS)
	}

	b.sendMarkdown(msg.Chat.ID, text)
}

func (b *Bot) send(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	b.api.Send(msg)
}

func (b *Bot) sendMarkdown(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}
