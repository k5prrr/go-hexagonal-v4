package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type TelegramConfig struct {
	Token   string `json:"Token"`
	Webhook bool   `json:"Webhook"`
}

type Telegram struct {
	config *TelegramConfig
}

func New(config *TelegramConfig) *Telegram {
	return &Telegram{
		config: config,
	}
}

func (telegram *Telegram) botUrl(command string) string {
	return fmt.Sprintf(
		"https://api.telegram.org/bot%s/%s",
		telegram.config.Token,
		command,
	)
}

func (telegram *Telegram) SendMassage(chatID int, message string, replyMarkup string) (string, error) {
	// Создаём данные сообщения
	messageMap := map[string]string{
		"chat_id": strconv.Itoa(chatID),
		"text":    message,
	}
	if replyMarkup != "" {
		messageMap["reply_markup"] = replyMarkup
	}
	//if (replyMarkup) {}

	// Перегоняем их в json
	messageJson, err := json.Marshal(messageMap)
	if err != nil {
		return "", err
	}

	response, err := http.Post(
		telegram.botUrl("sendMessage"),
		"application/json",
		bytes.NewBuffer(messageJson),
	)
	//fmt.Println(telegram.botUrl("sendMessage"))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Читаем тело ответа
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

func (telegram *Telegram) SendPhoto(chatID int, urlPhoto string, message string, replyMarkup string) (string, error) {
	// Создаём данные сообщения
	messageMap := map[string]string{
		"chat_id": strconv.Itoa(chatID),
		"photo":   urlPhoto,
	}
	if message != "" {
		messageMap["caption"] = message
	}
	if replyMarkup != "" {
		messageMap["reply_markup"] = replyMarkup
	}

	// Перегоняем их в json
	messageJson, err := json.Marshal(messageMap)
	if err != nil {
		return "", err
	}

	response, err := http.Post(
		telegram.botUrl("sendPhoto"),
		"application/json",
		bytes.NewBuffer(messageJson),
	)
	//fmt.Println(telegram.botUrl("sendMessage"))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Читаем тело ответа
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}

type InputMessage struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date    int    `json:"date"`
		Text    string `json:"text"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
	} `json:"message"`
}

func (inputMessage *InputMessage) New(inputBodyBytes *[]byte) error {
	return json.Unmarshal(*inputBodyBytes, inputMessage)
}

type SimpleInputMessage struct {
	ChatID      int      `json:"chat_id"`
	MessageID   int      `json:"message_id"`
	ClickButton bool     `json:"click_button"`
	Params      []string `json:"button_params"`
	Text        string   `json:"text"`
}

func (simpleInputMessage *SimpleInputMessage) FromInputMessage(inputMessage *InputMessage) {
	//fmt.Println(inputMessage)
	simpleInputMessage.ChatID = inputMessage.Message.Chat.ID
	simpleInputMessage.MessageID = inputMessage.Message.MessageID

	//simpleInputMessage.ClickButton = false
	//simpleInputMessage.Params = []
	simpleInputMessage.Text = inputMessage.Message.Text

	if inputMessage.Message.Contact.UserID != 0 {
		simpleInputMessage.ClickButton = true
		simpleInputMessage.Text = "contact"
		simpleInputMessage.Params = []string{
			strconv.Itoa(inputMessage.Message.Contact.UserID),
			inputMessage.Message.Contact.PhoneNumber,
		}
	}

}
func (simpleInputMessage *SimpleInputMessage) New(inputBodyBytes *[]byte) {
	inputMessage := &InputMessage{}
	inputMessage.New(inputBodyBytes)
	simpleInputMessage.FromInputMessage(inputMessage)
}
