package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
	"unicode"
)

func OnlyDigits(input string) string {
	var result []rune
	for _, char := range input {
		if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}
	return string(result)
}

func ClearConsole() {
	cmd := exec.Command("clear") // На Windows используйте "cls"
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func FormPhone(phone string) string {
	return fmt.Sprintf("+%s", phone)
}

func ReadJson(path string) (*map[string]interface{}, error) {

	// Проверка на существование файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	// Читаем файл
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// С Json в map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func IntValue(value interface{}, defaultValue int) int {
	switch val := value.(type) {
	case int:
		return val
	case float32, float64:
		return int(val.(float64))
	case int32, int64:
		return int(val.(int64))
	case string:
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal
		}
	}
	return defaultValue
}

func StringValue(value interface{}, defaultValue string) string {
	switch val := value.(type) {
	case string:
		return val
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", val)
	}
	return defaultValue
}

func FloatValue(value interface{}, defaultValue float64) float64 {
	if val, ok := value.(float64); ok {
		return val
	}
	return defaultValue
}
func TimeValue(value interface{}, defaultValue time.Time) time.Time {
	if val, ok := value.(time.Time); ok {
		return val
	}
	return defaultValue
}

// RandomInt generates a random integer between min and max (inclusive).
func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// RandomString generates a random string of specified length.
func RandomString(length int) string {
	if length <= 0 {
		length = 64
	}

	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	// -_.!*'()
	charLen := len(chars)
	result := make([]byte, length)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(charLen)]
	}

	return string(result)
}
