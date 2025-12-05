// Use
// env := env.New("")
// value := env.Get("key", "default")
package env

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	Path string
}

func New(envPath string) *Env {
	if envPath == "" {
		envPath = ".env"
	}

	env := Env{Path: envPath}
	env.init()
	return &env
}

func (e *Env) init() {
	file, err := os.Open(e.Path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		os.Setenv(key, value)
	}
}

func (e *Env) Get(key, defaultValue string) string {
	result, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return result
}
func (e *Env) Int(key string, defaultValue int64) int64 {
	resultStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	result, err := strconv.ParseInt(resultStr, 10, 64)
	if err != nil {
		return defaultValue
	}

	return result
}
