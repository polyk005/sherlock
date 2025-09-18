package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/polyk005/sherlock/pkg/channel"
	"github.com/polyk005/sherlock/pkg/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// Инициализация логгера
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Загрузка конфигурации
	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка загрузки конфигурации: %s", err.Error())
	}

	// Загрузка .env файла (если есть)
	if err := godotenv.Load(); err != nil {
		logrus.Info(".env файл не найден, используем конфиг")
	}

	// Создаем директорию для данных
	if err := os.MkdirAll("data", 0755); err != nil {
		logrus.Fatalf("Ошибка создания директории data: %s", err.Error())
	}

	// Инициализация базы данных
	dbPath := viper.GetString("database.path")
	db, err := repository.NewSQLiteDB(repository.SQLiteConfig{Path: dbPath})
	if err != nil {
		logrus.Fatalf("Ошибка инициализации БД: %s", err.Error())
	}
	defer db.Close()

	// Создание таблиц
	if err := repository.CreateSubscribersTable(db); err != nil {
		logrus.Fatalf("Ошибка создания таблиц: %s", err.Error())
	}

	// Получение конфигурации
	botToken := getConfigValue("BOT_TOKEN")
	channelID := viper.GetInt64("channel.id")
	adminChatID := viper.GetInt64("admin.chat_id")

	logrus.Infof("Запуск мониторинга канала %d", channelID)

	// Создание и запуск монитора
	monitor := channel.NewMonitor(botToken, channelID, adminChatID, db)
	if err := monitor.Start(); err != nil {
		logrus.Fatalf("Ошибка запуска монитора: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("database.path", "data/sherlock.db")
	viper.SetDefault("channel.id", -1001234567890)
	viper.SetDefault("admin.chat_id", 123456789)

	return viper.ReadInConfig()
}

func getConfigValue(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return viper.GetString(key)
}
