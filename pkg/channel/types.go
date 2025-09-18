package channel

import "time"

// Subscriber представляет подписчика канала
type Subscriber struct {
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	JoinedAt  time.Time `json:"joined_at"`
	LastCheck time.Time `json:"last_check"`
}

// Config представляет конфигурацию монитора
type Config struct {
	BotToken      string
	ChannelID     int64
	AdminChatID   int64
	CheckInterval time.Duration
}
