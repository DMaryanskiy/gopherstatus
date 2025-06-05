package bot

import (
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

// Mock database
type mockDB struct{}

func (m *mockDB) GetUserByTelegram(username string) (storage.User, error) {
	if username == "validuser" {
		return storage.User{ID: 1}, nil
	}
	return storage.User{}, errors.New("user not found")
}

func (m *mockDB) LatestResultsByID(limit int, userID uint) ([]storage.CheckResult, error) {
	return []storage.CheckResult{
		{ServiceName: "MockService", Online: true, ResponseMS: 123},
	}, nil
}

// Fake Telegram API
type fakeAPI struct {
	sent []tgbotapi.MessageConfig
}

func (f *fakeAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	msg := c.(tgbotapi.MessageConfig)
	f.sent = append(f.sent, msg)
	return tgbotapi.Message{}, nil
}
func (f *fakeAPI) StopReceivingUpdates() {}
func (f *fakeAPI) GetUpdatesChan(u tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	return make(chan tgbotapi.Update)
}

func TestNewBot_InvalidToken(t *testing.T) {
	_, err := NewBot("invalid-token", &storage.DB{})
	if err == nil {
		t.Error("expected error with invalid token")
	}
}

func TestHandleCommand_Help(t *testing.T) {
	b := &Bot{
		api: &fakeAPI{},
		db:  &mockDB{},
	}

	msg := &tgbotapi.Message{
		Text: "/help",
		Chat: &tgbotapi.Chat{ID: 123},
		From: &tgbotapi.User{UserName: "validuser"},
	}
	msg.Entities = []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}}

	b.handleCommand(msg)

	api := b.api.(*fakeAPI)
	if len(api.sent) != 1 {
		t.Fatalf("expected 1 message, got %d", len(api.sent))
	}
	if api.sent[0].Text == "" {
		t.Error("expected non-empty help message")
	}
}
