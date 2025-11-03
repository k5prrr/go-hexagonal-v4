package telegram

import (
	"app/internal/app/core/port"
	"app/pkg/telegram"
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

// Реагируем только на старт с верным кодом и на контакт

type TelegramWorker struct {
	tg      *telegram.Telegram
	useCase port.IUseCase
	cancel  context.CancelFunc
}

func NewTelegramWorker(token string, useCase port.IUseCase) *TelegramWorker {
	return &TelegramWorker{
		tg: telegram.New(&telegram.TelegramConfig{
			Token:   token,
			Webhook: false,
		}),
		useCase: useCase,
	}
}

func (t *TelegramWorker) Start(ctx context.Context, timeoutSec int) {
	ctx, t.cancel = context.WithCancel(ctx)
	interval := time.Duration(timeoutSec) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Telegram worker stopped")
			return
		case <-ticker.C:
			t.Step()
		}
	}
}
func (t *TelegramWorker) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
}

func (t *TelegramWorker) Step() {
	messages, err := t.tg.GetSimpleUpdates()
	if err != nil {
		log.Printf("failed to get updates: %v", err)
		return
	}

	for _, message := range messages {
		if strings.HasPrefix(message.Text, "/start ") {
			code := strings.TrimSpace(message.Text[len("/start "):])
			if code == "" {
				continue
			}

			log.Printf("Received /start with code: %s from chatID=%d", code, message.ChatID)

			ctx := context.Context(context.Background())
			err := t.useCase.AddChatIdByCodeCheckPhone(ctx, code, message.ChatID)
			if err != nil {
				log.Printf("failed AddChatIdByCodeCheckPhone: %v", err)
			}
			continue
		}
		if message.Text == "contact" &&
			message.ClickButton &&
			strconv.Itoa(message.ChatID) == message.Params[0] {

			log.Printf("Received contact from chatID=%d: UserID=%s, Phone=%s",
				message.ChatID,
				message.Params[0],
				message.Params[1],
			)

			ctx := context.Context(context.Background())
			err := t.useCase.AddPhoneByChatId(ctx, message.ChatID, message.Params[1])
			if err != nil {
				log.Printf("failed AddPhoneByChatId: %v", err)
			}
			continue
		}
		// Тут если он уже есть в системе и номера нет, то отправлять его на отправку номера
	}
}
func (t *TelegramWorker) acceptCode(chatID int, code string) {

}
func (t *TelegramWorker) acceptContact(chatID int, phone string) {

}
