package telegramW

import (
	"app/internal/app/core/port"
	"app/pkg/telegram"
	"context"
	"fmt"
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

func New(tg      *telegram.Telegram, useCase port.IUseCase) *TelegramWorker {
	return &TelegramWorker{
		tg: tg,
		useCase: useCase,
	}
}

func (t *TelegramWorker) Run(ctx context.Context, timeoutSec int) {
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
		textList := strings.Split(message.Text, "/start ")

		if textList[0] == "" { // Сам старт

			if textList[1] == "registration" {
				// Отправляем запрос номера
				_, err = t.tg.SendMessage(
					message.ChatID,
					"Нажмите на кнопку \"Отправить номер\" на дополнительной клавиатуре",
					"{\"keyboard\":[[{\"text\":\"Отправить номер\",\"request_contact\":true}]],\"one_time_keyboard\":true}",
				)
				if err != nil {
					log.Printf("failed to send message: %v", err)
				}
				continue
			}
		}

		if message.Text == "contact" &&
			message.ClickButton &&
			strconv.FormatInt(message.ChatID, 10) == message.Params[0] {
			// Получаем ответ его контакт

			log.Printf("Received contact from chatID=%d: UserID=%s, Phone=%s",
				message.ChatID,
				message.Params[0],
				message.Params[1],
			)

			_, err = t.tg.SendMessage(
				message.ChatID,
				fmt.Sprintf("Номер принят!\n\nВернитесь на сайт и войдите с номером\n%s", message.Params[1]),
				"{\"hide_keyboard\": true}",
			)
			if err != nil {
				log.Printf("failed to send message: %v", err)
			}

			phone, err := strconv.ParseInt(message.Params[1], 10, 64)
			if err != nil {
				log.Printf("failed to parse phone: %v", err)
			}
			ctx := context.Context(context.Background())
			_, err = t.useCase.CreateUser(ctx, message.ChatID, phone)
			if err != nil {
				log.Printf("failed CreateUser: %v", err)
			}
			continue
		}

		/*if strings.HasPrefix(message.Text, "/start ") {
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
		}*/

		// Тут если он уже есть в системе и номера нет, то отправлять его на отправку номера
	}
}
