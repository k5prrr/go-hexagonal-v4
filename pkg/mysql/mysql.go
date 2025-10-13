package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	//	"regexp"
	"time"
	//"log"

	_ "github.com/go-sql-driver/mysql"
)

type SqlConfig struct {
	Host     string
	Port     string
	Login    string
	Password string
	Name     string
	Charset  string

	MaxIdleConns          int
	MaxOpenConns          int
	ConnMaxLifetimeMinute time.Duration
}

type Mysql struct {
	DB     *sql.DB
	config *SqlConfig
}

func New(config *SqlConfig) (*Mysql, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s",
			config.Login,
			config.Password,
			config.Host,
			config.Port,
			config.Name,
			config.Charset,
		),
	)

	if err != nil {
		return nil, err
	}

	// Устанавливаем максимальное количество простаивающих соединений
	db.SetMaxIdleConns(config.MaxIdleConns)

	// Устанавливаем максимальное количество открытых соединений
	db.SetMaxOpenConns(config.MaxOpenConns)

	// Устанавливаем максимальное время жизни соединения
	db.SetConnMaxLifetime(time.Minute * config.ConnMaxLifetimeMinute) // time.Duration(

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Mysql{
		DB:     db,
		config: config,
	}, nil
}

func (mysql *Mysql) Stop() error {
	return mysql.DB.Close()
}

// Функция, которая просто делает sql запрос и не возвращает ответ
func (mysql *Mysql) Query(request string) error {
	re := regexp.MustCompile(`;\s*`)
	requests := re.Split(request, -1)

	for _, line := range requests {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		_, err := mysql.DB.Exec(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// Функция, которая выбирает данные и возвращает в виде массива с мапами (либо слайса) в переданный ей интерфейс
func (mysql *Mysql) SelectQuery(request string) ([]map[string]interface{}, error) {
	//fmt.Println(request)

	// Сам запрос
	rows, err := mysql.DB.Query(request)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Название столбцов в виде массива
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Создаем массив для хранения результатов
	var results []map[string]interface{}

	// Проходим по всем строкам результата запроса
	for rows.Next() {
		// Создаем массив для хранения значений столбцов
		values := make([]interface{}, len(columns))
		// Создаем массив для хранения указателей на значения столбцов
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		// Сканируем строку в массив значений
		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		// Создаем карту для хранения значений столбцов
		entry := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				entry[colName] = string(b)
			} else {
				entry[colName] = val
			}
		}

		// Добавляем карту в массив результатов
		results = append(results, entry)
	}

	// Проверяем, есть ли ошибки после цикла
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возвращаем результаты
	return results, nil
}

// Тоже самое, что и с select, но возвращает только 1 элемент в виде карты
func (mysql *Mysql) SelectOne(request string) (map[string]interface{}, error) {
	// Вызываем SelectQuery с добавлением LIMIT 1 к запросу
	result, err := mysql.SelectQuery(fmt.Sprintf("%s LIMIT 1", request))
	if err != nil {
		return nil, err
	}

	// Проверяем, есть ли результаты
	if len(result) > 0 {
		// Возвращаем первый элемент из результатов
		return result[0], nil
	}

	// Если результатов нет, возвращаем nil
	return nil, nil
}

// Добавляет данные в таблицу
func (mysql *Mysql) Entry(tableName string, data map[string]interface{}) error {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		values = append(values, mysql.ClearText(fmt.Sprintf("%v", val)))
	}

	placeholders := strings.Repeat("?, ", len(values))
	placeholders = placeholders[:len(placeholders)-2] // remove trailing ", "

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), placeholders)
	_, err := mysql.DB.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("query: %s, error: %w", query, err)
	}
	return nil
}

// Добавляет множество данных в таблицу
func (mysql *Mysql) EntryMany(tableName string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return nil
	}

	columns := make([]string, 0, len(data[0]))
	for col := range data[0] {
		columns = append(columns, col)
	}

	placeholders := strings.Repeat("(?), ", len(columns))
	placeholders = placeholders[:len(placeholders)-2] // remove trailing ", "

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", mysql.ClearName(tableName), strings.Join(columns, ", "), placeholders)

	tx, err := mysql.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range data {
		values := make([]interface{}, 0, len(item))
		for _, val := range item {
			values = append(values, mysql.ClearText(fmt.Sprintf("%v", val)))
		}
		_, err := stmt.Exec(values...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Меняет елемент по его ид
func (mysql *Mysql) Update(tableName string, id int, data map[string]interface{}) error {
	// Проверка на пустую map
	if len(data) == 0 {
		return fmt.Errorf("data map is empty, nothing to update")
	}

	setStatements := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for col, val := range data {
		setStatements = append(setStatements, fmt.Sprintf("%s = ?", mysql.ClearName(col)))
		values = append(values, fmt.Sprintf("%v", val)) // mysql.ClearText(fmt.Sprintf("%v", val))
	}
	values = append(values, id)

	queryTop := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", mysql.ClearName(tableName), strings.Join(setStatements, ", "))
	//query := fmt.Sprintf(queryTop, values...)
	_, err := mysql.DB.Exec(queryTop, values...)
	if err != nil {
		return fmt.Errorf("query: %s, error: %w", queryTop, err, values)
	}
	return nil
}

// Удаляет элемент по его ид
func (mysql *Mysql) Delete(tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", mysql.ClearName(tableName))
	_, err := mysql.DB.Exec(query, id)
	return err
}

// Очищает имя подготавливая его к базе данных
func (mysql *Mysql) ClearName(text string) string {
	text = strings.TrimSpace(text)
	text = regexp.MustCompile("[^a-zA-Z0-9_]").ReplaceAllString(text, "")
	text = strings.TrimSpace(text)
	if len(text) > 64 { // MySQL имена столбцов не могут быть длиннее 64 символов
		text = text[:64]
	}
	return text
}

// Очищает текст подготавливая его к базе данных
func (mysql *Mysql) ClearText(text string) string {
	if json.Valid([]byte(text)) {
		// Если текст является валидным JSON, возвращаем его без изменений
		return text
	}

	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\r", "\\r")  // Экранируем символ перевода каретки
	text = strings.ReplaceAll(text, "\"", "\\\"") // Экранируем двойные кавычки
	text = strings.ReplaceAll(text, "'", "\\'")   // Экранируем одинарные кавычки
	text = strings.ReplaceAll(text, "`", "\\`")   // Экранируем обратные кавычки
	text = strings.ReplaceAll(text, "\n", "\\n")  // Экранируем символ новой строки
	return text
}

// Показ результата
func (mysql *Mysql) ResultToString(data []map[string]interface{}) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Error ResultToString data: %v \n err: %v", data, err)
	}

	// Преобразование []byte в строку
	jsonString := string(jsonData)

	return jsonString, nil
}
