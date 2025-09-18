package channel

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func (m *Monitor) notifyNewSubscriber(subscriber Subscriber) {
	message := fmt.Sprintf("🎉 *Новый подписчик!*\n\n"+
		"👤 Имя: %s %s\n"+
		"📎 Username: @%s\n"+
		"🆔 User ID: `%d`\n"+
		"⏰ Время: %s",
		subscriber.FirstName, subscriber.LastName,
		subscriber.Username,
		subscriber.UserID,
		time.Now().Format("15:04:05"),
	)

	msg := tgbotapi.NewMessage(m.config.AdminChatID, message)
	msg.ParseMode = "Markdown"
	m.bot.Send(msg)
}

func (m *Monitor) notifyUnsubscribed(subscriber Subscriber) {
	duration := time.Since(subscriber.JoinedAt).Round(time.Second)

	// Форматируем длительность по-русски
	durationStr := formatDurationRu(duration)

	message := fmt.Sprintf("😢 *Кто-то отписался*\n\n"+
		"👤 Имя: %s %s\n"+
		"📎 Username: @%s\n"+
		"🆔 User ID: `%d`\n"+
		"⏰ Подписан: %s\n"+
		"⏱️ Был с нами: %s",
		subscriber.FirstName, subscriber.LastName,
		subscriber.Username,
		subscriber.UserID,
		subscriber.JoinedAt.Format("02.01.2006 15:04:05"),
		durationStr,
	)

	msg := tgbotapi.NewMessage(m.config.AdminChatID, message)
	msg.ParseMode = "Markdown"

	// Пытаемся отправить уведомление
	if _, err := m.bot.Send(msg); err != nil {
		logrus.Errorf("Ошибка отправки уведомления об отписке: %v", err)
		// Сохраняем в лог для последующей отправки
		m.savePendingNotification(message)
	}
}

// formatDurationRu форматирует длительность по-русски
func formatDurationRu(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string

	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", days, pluralizeRu(days, "день", "дня", "дней")))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", hours, pluralizeRu(hours, "час", "часа", "часов")))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d %s", minutes, pluralizeRu(minutes, "минуту", "минуты", "минут")))
	}

	if len(parts) == 0 {
		return "менее минуты"
	}

	return strings.Join(parts, " ")
}

// pluralizeRu склонение русских существительных
func pluralizeRu(n int, singular, few, many string) string {
	n = n % 100
	if n > 10 && n < 20 {
		return many
	}

	n = n % 10
	if n == 1 {
		return singular
	}
	if n >= 2 && n <= 4 {
		return few
	}
	return many
}

// savePendingNotification сохраняет уведомление для последующей отправки
func (m *Monitor) savePendingNotification(message string) {
	// Можно сохранить в файл или временную базу данных
	logrus.Warnf("Неотправленное уведомление: %s", message)
}
