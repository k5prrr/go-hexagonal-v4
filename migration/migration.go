package migration

import (
	"context"
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// migrationsFS содержит все SQL-файлы из директории migrations/ в едином
// embed.FS образе. Это гарантирует, что миграции "вшиты" в бинарник
// и не зависят от внешней файловой системы - критично для контейнеров и CI/CD.
// Обратите внимание: встраиваем **только** миграции (не down-файлы, не .md и др.),
// фильтрация происходит по расширению `.up.sql`.
//
//go:embed migrations/*.sql
var migrationsFS embed.FS

// migration представляет *одну* миграцию как объект в памяти.
// Поля:
//   - Version: числовая версия (например, 1, 2, 15), извлекается из имени файла (001_, 015_ и т.д.)
//   - Name: точное имя файла (например, "001_create_users.up.sql")
//   - Content: полный текст SQL-скрипта без обработки/интерпретации
//
// Не экспортируемый тип - используется только внутри пакета, что уменьшает surface API.
type migration struct {
	Version int64
	Name    string
	Content string
}

// Migration - основной тип, инкапсулирующий логику применения миграций.
// Хранит пул соединений к БД и реализует идемпотентное применение миграций
// в строгом порядке возрастания версий.
//
// Особенности:
//   - Использует только std и github.com/jackc/pgx/v5 (без сторонних миграционных фреймворков)
//   - Каждая миграция выполняется в отдельной транзакции
//   - После успешного выполнения SQL миграции её версия атомарно записывается в schema_migrations
//   - Поддерживает идемпотентность: повторный вызов Run() ничего не делает, если все миграции применены
type Migration struct {
	pool *pgxpool.Pool
}

// New создаёт новый экземпляр Migration с заданным пулом соединений.
// Пул должен быть уже инициализирован и работоспособен.
// Рекомендуется передавать пул с ограниченным количеством соединений (например, max 1–2 для миграций),
// так как миграции - не высоконагруженная операция, и избыточные соединения не нужны.
func New(pool *pgxpool.Pool) *Migration {
	return &Migration{pool: pool}
}

// Run применяет все **отложенные (pending)** миграции в порядке возрастания версий.
// Логика:
//   1. Загружает список уже применённых версий из таблицы schema_migrations
//   2. Сканирует встроенные (embed) SQL-файлы, парсит версии из имён
//   3. Исключает уже применённые миграции
//   4. Сортирует оставшиеся по Version
//   5. Выполняет каждую в отдельной транзакции (SQL + INSERT INTO schema_migrations)
//
// Гарантии:
//   - Атомарность: либо миграция полностью применена (SQL + запись версии), либо не применена совсем
//   - Порядок: строгая сортировка по Version исключает race-условии при параллельном запуске
//   - Идемпотентность: повторный вызов Run() безопасен - он пропустит уже применённые миграции
//
// Примечание: таблица `schema_migrations` должна существовать до первого вызова Run().
// Обычно её создают отдельной инициализационной миграцией (например, 000_init_schema_migrations.up.sql)
// или вручную: CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT PRIMARY KEY);
func (m *Migration) Run(ctx context.Context) error {
	// Шаг 1: получаем список версий, уже применённых в БД
	applied, err := m.appliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Шаг 2: загружаем и фильтруем миграции, которые ещё не применены
	pending, err := m.loadPendingMigrations(applied)
	if err != nil {
		return fmt.Errorf("failed to load pending migrations: %w", err)
	}

	// Шаг 3: применяем каждую миграцию в отдельной транзакции
	for _, mig := range pending {
		// Прогресс: выводим в stderr/stdout для логирования в CI/CD или при запуске вручную
		// Формат: "Applying migration 001: 001_create_users.up.sql"
		// Ведущие нули - для визуального выравнивания в логах.
		fmt.Printf("Applying migration %03d: %s\n", mig.Version, mig.Name)

		// Начинаем транзакцию. Если миграция "тяжёлая", убедитесь, что в контексте
		// нет слишком короткого таймаута (например, context.WithTimeout).
		tx, err := m.pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", mig.Name, err)
		}
		//️ Безусловный defer tx.Rollback() - классический паттерн:
		// - если будет вызван tx.Commit() - Rollback ничего не сделает (транзакция уже закрыта)
		// - если ошибка или паника - транзакция гарантированно откатится
		defer tx.Rollback(ctx)

		// Выполняем SQL миграции "как есть" - без подстановок, без шаблонизации.
		// Это упрощает отладку и делает миграции прозрачными.
		_, err = tx.Exec(ctx, mig.Content)
		if err != nil {
			return fmt.Errorf("migration %s failed (SQL execution): %w", mig.Name, err)
		}

		// Фиксируем факт применения миграции: записываем её версию в мета-таблицу.
		// Это делается в той же транзакции - обеспечивая атомарность:
		// либо SQL выполнится И версия запишется, либо ничего не произойдёт.
		_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", mig.Version)
		if err != nil {
			// Например, UNIQUE VIOLATION, если где-то race condition (но в нашем случае - маловероятно)
			return fmt.Errorf("failed to record migration %d in schema_migrations: %w", mig.Version, err)
		}

		// Коммитим транзакцию. Если ошибка - миграция не применена, и при следующем запуске
		// её попробуют применить снова (идемпотентность соблюдена).
		err = tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("failed to commit transaction for migration %d: %w", mig.Version, err)
		}
	}

	// Логика завершена успешно: все pending миграции применены.
	return nil
}

// appliedMigrations запрашивает список уже применённых версий из БД.
// Возвращает map[int64]bool для O(1)-проверки "была ли миграция применена".
// Предполагает, что таблица `schema_migrations` существует и содержит колонку `version BIGINT`.
// Ошибки запроса (например, отсутствие таблицы) пробрасываются наверх - это не "нормальная" ситуация.
func (m *Migration) appliedMigrations(ctx context.Context) (map[int64]bool, error) {
	rows, err := m.pool.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err // например, relation "schema_migrations" does not exist
	}
	defer rows.Close() // освобождаем ресурсы - важно при повторных вызовах

	applied := make(map[int64]bool)
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("failed to scan version from schema_migrations: %w", err)
		}
		applied[v] = true
	}
	// Проверяем, не было ли ошибки в итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating schema_migrations rows: %w", err)
	}
	return applied, nil
}

// loadPendingMigrations сканирует встроенные файлы миграций и отбирает те,
// которых ещё нет в `applied`.
//
// Соглашение об именовании файлов:
//   - Формат: `<N>_<описание>.up.sql`
//   - Примеры: `001_init.up.sql`, `015_add_index.up.sql`
//   - N - номер версии, допускаются ведущие нули (парсится как int64)
//   - Обязательно расширение `.up.sql` (`.down.sql` игнорируются - их нет в embed)
//
// Алгоритм:
//   1. Читаем список файлов из embed.FS
//   2. Фильтруем по расширению `.up.sql`
//   3. Извлекаем номер версии из префикса имени файла (до первого '_')
//   4. Пропускаем, если applied[version] == true
//   5. Читаем содержимое файла и добавляем в список
//   6. Сортируем по Version - даже если файловая система вернёт их в "неправильном" порядке
//
// Важно: не поддерживает отрицательные или нечисловые префиксы.
// Не поддерживает "дырявую" нумерацию (001, 003, 002) - но сортировка это исправляет.
func (m *Migration) loadPendingMigrations(applied map[int64]bool) ([]migration, error) {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		// Может возникнуть, если путь "migrations" не существует в embed (но go:embed должен был этого не допустить)
		return nil, fmt.Errorf("failed to read embedded migrations directory: %w", err)
	}

	var pending []migration
	for _, entry := range entries {
		name := entry.Name()

		// Поддерживаем ТОЛЬКО .up.sql - это упрощает логику и исключает случайное применение down-миграций
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		// Разбиваем имя файла на части: "001_create_users.up.sql" → ["001", "create_users.up.sql"]
		parts := strings.SplitN(name, "_", 2)
		if len(parts) < 1 {
			return nil, fmt.Errorf("invalid migration filename (missing '_'): %s", name)
		}

		// Извлекаем числовую часть: удаляем ведущие нули (но "000" → "0")
		vStr := strings.TrimLeft(parts[0], "0")
		if vStr == "" {
			vStr = "0" // случай "000_.up.sql" → версия 0
		}

		// Парсим как int64 - позволяет иметь до 9 квинтиллионов миграций
		v, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid version prefix in %s (expected digits only): %w", name, err)
		}

		// Пропускаем, если уже применена
		if applied[v] {
			continue
		}

		// Читаем содержимое SQL-файла из embed.FS
		content, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded migration %s: %w", name, err)
		}

		pending = append(pending, migration{
			Version: v,
			Name:    name,
			Content: string(content),
		})
	}

	// Обязательная сортировка - даже если файлы в embed идут в правильном порядке,
	// это не гарантируется спецификацией (особенно для разных ОС/файловых систем)
	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})

	return pending, nil
}