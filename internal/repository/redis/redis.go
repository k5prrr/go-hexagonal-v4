// go get -u github.com/go-redis/redis/v8
package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:45013",
		Password: "",
		DB:       0,
	})

	// Запись
	err := rdb.Set(ctx, "key1", "value1", 0).Err()
	if err != nil {
		log.Fatalf("Ошибка при записи в Redis: %v", err)
	}
	fmt.Println("Данные записаны в Redis")

	// Чтение
	val, err := rdb.Get(ctx, "key1").Result()
	if err != nil {
		log.Fatalf("Ошибка при чтении из Redis: %v", err)
	}
	fmt.Printf("Полученные данные: %s\n", val)

	// Изменение
	err = rdb.Set(ctx, "key1", "new_value1", 0).Err()
	if err != nil {
		log.Fatalf("Ошибка при изменении данных в Redis: %v", err)
	}
	fmt.Println("Данные изменены в Redis")

	// Удаление данных из Redis
	err = rdb.Del(ctx, "key1").Err()
	if err != nil {
		log.Fatalf("Ошибка при удалении данных из Redis: %v", err)
	}
	fmt.Println("Данные удалены из Redis")

	// ---

	// Проверка измененных данных
	newVal, err := rdb.Get(ctx, "key1").Result()
	if err != nil {
		log.Fatalf("Ошибка при чтении измененных данных из Redis: %v", err)
	}
	fmt.Printf("Измененные данные: %s\n", newVal)

	// Попытка чтения удаленных данных
	deletedVal, err := rdb.Get(ctx, "key1").Result()
	if err == redis.Nil {
		fmt.Println("Ключ 'key1' не существует в Redis")
	} else if err != nil {
		log.Fatalf("Ошибка при чтении удаленных данных из Redis: %v", err)
	} else {
		fmt.Printf("Удаленные данные: %s\n", deletedVal)
	}
}
