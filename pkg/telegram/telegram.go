package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
https://api.telegram.org/bot<TOKEN>/getUpdates
curl "https://api.telegram.org/bot<TOKEN>/deleteWebhook"
*/
var (
	offsetFilePath = "data/telegramOffset.txt"
)

type TelegramConfig struct {
	Token   string `json:"Token"`
	Webhook bool   `json:"Webhook"`
}

type Telegram struct {
	config *TelegramConfig
	offset int64
	mu     sync.RWMutex
}

func New(config *TelegramConfig) *Telegram {
	t := &Telegram{
		config: config,
	}
	t.loadOffset()
	return t
}
func (t *Telegram) loadOffset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	data, err := os.ReadFile(offsetFilePath)
	if err != nil {
		// Файл не существует или ошибка — начинаем с 0
		t.offset = 0
		return
	}

	str := string(data)
	str = strings.TrimSpace(str)
	if str == "" {
		t.offset = 0
		return
	}

	parsed, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		t.offset = 0
		return
	}

	t.offset = parsed
}

func (t *Telegram) botUrl(command string) string {
	return fmt.Sprintf(
		"https://api.telegram.org/bot%s/%s",
		t.config.Token,
		command,
	)
}

func (t *Telegram) SendMassage(chatID int, message string, replyMarkup string) (string, error) {
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
		t.botUrl("sendMessage"),
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

func (t *Telegram) SendPhoto(chatID int, urlPhoto string, message string, replyMarkup string) (string, error) {
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
		t.botUrl("sendPhoto"),
		"application/json",
		bytes.NewBuffer(messageJson),
	)
	//fmt.Println(t.botUrl("sendMessage"))
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
func (t *Telegram) getOffset() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.offset
}
func (t *Telegram) setOffset(offset int64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.offset = offset

	return os.WriteFile(offsetFilePath, []byte(strconv.FormatInt(offset, 10)), 0644)
}
func (t *Telegram) GetUpdates() ([]InputMessage, error) {
	url := t.botUrl(fmt.Sprintf("getUpdates?offset=%d&timeout=30", t.getOffset()+1))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read updates: %w", err)
	}

	var result InputMessages
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updates: %w", err)
	}

	if !result.Ok {
		return nil, fmt.Errorf("API returned error: %s", string(body))
	}

	lenResults := len(result.Result)
	if lenResults != 0 {
		if err := t.setOffset(int64(result.Result[lenResults-1].UpdateID)); err != nil {
			return nil, err
		}
	}

	return result.Result, nil
}

type InputMessages struct {
	Ok     bool           `json:"ok"`
	Result []InputMessage `json:"result"`
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
