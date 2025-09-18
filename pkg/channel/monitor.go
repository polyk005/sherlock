package channel

import (
	"database/sql"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Monitor struct {
	bot    *tgbotapi.BotAPI
	db     *sql.DB
	config Config
}

func NewMonitor(botToken string, channelID, adminChatID int64, db *sql.DB) *Monitor {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logrus.Fatalf("Ошибка создания бота: %v", err)
	}

	return &Monitor{
		bot: bot,
		db:  db,
		config: Config{
			BotToken:      botToken,
			ChannelID:     channelID,
			AdminChatID:   adminChatID,
			CheckInterval: 2 * time.Minute,
		},
	}
}

func (m *Monitor) Start() error {
	logrus.Info("Мониторинг запущен")

	go m.periodicCheck()
	return m.handleRealTimeEvents()
}

func (m *Monitor) periodicCheck() {
	ticker := time.NewTicker(m.config.CheckInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := m.checkSubscribers(); err != nil {
			logrus.Errorf("Ошибка проверки подписчиков: %v", err)
		}

		// Периодическая проверка доступности бота
		if err := m.checkBotAvailability(); err != nil {
			logrus.Errorf("Бот недоступен: %v", err)
		}
	}
}

// checkBotAvailability проверяет, что бот может отправлять сообщения
func (m *Monitor) checkBotAvailability() error {
	_, err := m.bot.GetMe()
	return err
}

// checkSubscribers проверяет изменения в подписчиках
func (m *Monitor) checkSubscribers() error {
	// Получаем текущее количество участников
	memberCount, err := m.bot.GetChatMembersCount(tgbotapi.ChatMemberCountConfig{
		ChatConfig: tgbotapi.ChatConfig{ChatID: m.config.ChannelID},
	})
	if err != nil {
		return fmt.Errorf("ошибка получения количества участников: %v", err)
	}

	logrus.Infof("Текущее количество подписчиков: %d", memberCount)
	return nil
}

// handleChatMemberUpdate обрабатывает изменения статуса участников
func (m *Monitor) handleChatMemberUpdate(update *tgbotapi.ChatMemberUpdated) {
	if update.Chat.ID == m.config.ChannelID {
		user := update.NewChatMember.User
		oldStatus := update.OldChatMember.Status
		newStatus := update.NewChatMember.Status

		logrus.Infof("Изменение статуса: %s -> %s для пользователя %d",
			oldStatus, newStatus, user.ID)

		// Сохраняем информацию о пользователе перед обработкой
		subscriber := Subscriber{
			UserID:    user.ID,
			Username:  user.UserName,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			JoinedAt:  time.Now(), // Временно, потом заменим из базы
			LastCheck: time.Now(),
		}

		if oldStatus == "left" && newStatus == "member" {
			// Новый подписчик
			subscriber.JoinedAt = time.Now()
			if err := m.saveSubscriber(subscriber); err != nil {
				logrus.Errorf("Ошибка сохранения подписчика: %v", err)
				return
			}
			m.notifyNewSubscriber(subscriber)

		} else if oldStatus == "member" && newStatus == "left" {
			// Кто-то отписался - сначала получаем реальные данные из базы
			if existingSub, err := m.getSubscriberFromDB(user.ID); err == nil {
				// Используем реальные данные из базы
				subscriber = existingSub
			}

			// Удаляем из базы и уведомляем
			if err := m.deleteSubscriber(user.ID); err != nil {
				logrus.Errorf("Ошибка удаления подписчика: %v", err)
			}
			m.notifyUnsubscribed(subscriber)
		}
	}
}

// getSubscriberFromDB получает информацию о подписчике из базы данных
func (m *Monitor) getSubscriberFromDB(userID int64) (Subscriber, error) {
	var subscriber Subscriber
	var joinedAtStr string

	err := m.db.QueryRow(`
        SELECT user_id, username, first_name, last_name, joined_at 
        FROM subscribers WHERE user_id = ?
    `, userID).Scan(
		&subscriber.UserID,
		&subscriber.Username,
		&subscriber.FirstName,
		&subscriber.LastName,
		&joinedAtStr,
	)

	if err != nil {
		return subscriber, err
	}

	// Парсим время из строки
	subscriber.JoinedAt, _ = time.Parse(time.RFC3339, joinedAtStr)
	return subscriber, nil
}

func (m *Monitor) handleRealTimeEvents() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"chat_member", "message", "channel_post", "chat_join_request"}

	updates := m.bot.GetUpdatesChan(u)
	logrus.Info("Ожидание событий в канале...")

	for update := range updates {
		// Запрос на вступление в канал
		if update.ChatJoinRequest != nil {
			m.handleChatJoinRequest(update.ChatJoinRequest)
		}
		// Изменение статуса участника
		if update.ChatMember != nil {
			m.handleChatMemberUpdate(update.ChatMember)
		}
	}
	return nil
}

// Обработка запросов на вступление
func (m *Monitor) handleChatJoinRequest(request *tgbotapi.ChatJoinRequest) {
	user := request.From
	subscriber := Subscriber{
		UserID:    user.ID,
		Username:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		JoinedAt:  time.Now(),
		LastCheck: time.Now(),
	}

	if err := m.saveSubscriber(subscriber); err != nil {
		logrus.Errorf("Ошибка сохранения: %v", err)
	} else {
		logrus.Infof("Сохранен пользователь из запроса: %s", user.FirstName)
	}
}
