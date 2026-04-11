package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Драйвер Postgres
	"gopkg.in/telebot.v3"
)

func main() {
	// 1. Получаем настройки из окружения
	dsn := os.Getenv("DB_DSN")
	token := os.Getenv("BOT_TOKEN")

	// 2. Логика "Retry" для подключения к БД (сервер БД грузится дольше бота)
	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil && db.Ping() == nil {
			fmt.Println("✅ Успешное подключение к БД!")
			break
		}
		fmt.Printf("⚠️ Ожидание БД (попытка %d/5)...\n", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil || db.Ping() != nil {
		log.Fatal("❌ Не удалось подключиться к БД:", err)
	}

	// 3. Настройка бота
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, _ := telebot.NewBot(pref)

	// 4. Пример взаимодействия (Клиент -> Сервер БД)
	b.Handle("/start", func(c telebot.Context) error {
		// Записываем юзера в базу (упрощенно)
		_, err := db.Exec("INSERT INTO users (tg_id) VALUES ($1) ON CONFLICT DO NOTHING", c.Sender().ID)
		if err != nil {
			return c.Send("Ошибка при сохранении в БД")
		}
		return c.Send("Привет! Ты добавлен в мою базу данных.")
	})

	fmt.Println("🤖 Бот запущен...")
	b.Start()
}
