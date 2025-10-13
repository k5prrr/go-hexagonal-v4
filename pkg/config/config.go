/*
 * Use
 * conf := config.New("")
 * conf.Int("name/name")
 */
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Path    string
	Json    map[string]interface{}
	Strings map[string]string
	Bools   map[string]bool
	Ints    map[string]int
	Floats  map[string]float64
}

func New(configPath string) *Config {
	if configPath == "" {
		configPath = "configs/configs.json"
	}
	return &Config{
		Path:    configPath,
		Strings: make(map[string]string),
		Bools:   make(map[string]bool),
		Ints:    make(map[string]int),
		Floats:  make(map[string]float64),
	}
}

func (c *Config) ClearCache() {
	c.Strings = make(map[string]string)
	c.Bools = make(map[string]bool)
	c.Ints = make(map[string]int)
	c.Floats = make(map[string]float64)
}

func (c *Config) Update() error {
	text, err := os.ReadFile(c.Path)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	err = json.Unmarshal(text, &c.Json)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}

	c.ClearCache()
	return nil
}

func (c *Config) Current(name string) (interface{}, error) {
	if c.Json == nil {
		if err := c.Update(); err != nil {
			return nil, err
		}
	}

	keys := strings.Split(name, "/")
	var current interface{} = c.Json

	for _, key := range keys {
		if currentMap, ok := current.(map[string]interface{}); ok {
			if val, exists := currentMap[key]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("key '%s' not found in config", key)
			}
		} else {
			return nil, fmt.Errorf("key '%s' is not a map", key)
		}
	}

	return current, nil
}

func (c *Config) String(name string) (string, error) {
	value, exists := c.Strings[name]
	if exists {
		return value, nil
	}

	current, err := c.Current(name)
	if err != nil {
		return "", err
	}

	switch v := current.(type) {
	case string:
		c.Strings[name] = v
		return v, nil
	case int:
		c.Strings[name] = strconv.Itoa(v)
		return c.Strings[name], nil
	case bool:
		if v {
			c.Strings[name] = "true"
			return c.Strings[name], nil
		} else {
			c.Strings[name] = "false"
			return c.Strings[name], nil
		}
	case float64:
		c.Strings[name] = strconv.FormatFloat(v, 'f', -1, 64)
		return c.Strings[name], nil
	default:
		return "", fmt.Errorf("unsupported type '%s' for key '%s'", reflect.TypeOf(v).String(), name)
	}
}

func (c *Config) Bool(name string) (bool, error) {
	value, exists := c.Bools[name]
	if exists {
		return value, nil
	}

	current, err := c.Current(name)
	if err != nil {
		return false, err
	}

	switch v := current.(type) {
	case bool:
		c.Bools[name] = v
		return v, nil
	case int:
		c.Bools[name] = v != 0
		return c.Bools[name], nil
	case string:
		boolValue, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("cannot convert '%s' to bool", v)
		}
		c.Bools[name] = boolValue
		return boolValue, nil
	case float64:
		c.Bools[name] = v != 0
		return c.Bools[name], nil
	default:
		return false, fmt.Errorf("unsupported type '%s' for key '%s'", reflect.TypeOf(v).String(), name)
	}
}

func (c *Config) Int(name string) (int, error) {
	value, exists := c.Ints[name]
	if exists {
		return value, nil
	}

	current, err := c.Current(name)
	if err != nil {
		return 0, err
	}

	switch v := current.(type) {
	case int:
		c.Ints[name] = v
		return v, nil
	case float64:
		c.Ints[name] = int(v)
		return c.Ints[name], nil
	case string:
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("cannot convert '%s' to int", v)
		}
		c.Ints[name] = intValue
		return c.Ints[name], nil
	case bool:
		if v {
			c.Ints[name] = 1
		} else {
			c.Ints[name] = 0
		}
		return c.Ints[name], nil
	default:
		return 0, fmt.Errorf("unsupported type '%s' for key '%s'", reflect.TypeOf(v).String(), name)
	}
}

func (c *Config) Float(name string) (float64, error) {
	value, exists := c.Floats[name]
	if exists {
		return value, nil
	}

	current, err := c.Current(name)
	if err != nil {
		return 0, err
	}

	switch v := current.(type) {
	case float64:
		c.Floats[name] = v
		return v, nil
	case int:
		c.Floats[name] = float64(v)
		return c.Floats[name], nil
	case string:
		floatValue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert '%s' to float64", v)
		}
		c.Floats[name] = floatValue
		return floatValue, nil
	case bool:
		if v {
			return 1.0, nil
		} else {
			return 0.0, nil
		}
	default:
		return 0, fmt.Errorf("unsupported type '%s' for key '%s'", reflect.TypeOf(v).String(), name)
	}
}
