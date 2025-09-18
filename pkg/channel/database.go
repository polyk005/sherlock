package channel

import (
	"time"
)

func (m *Monitor) saveSubscriber(subscriber Subscriber) error {
	query := `
		INSERT OR REPLACE INTO subscribers 
		(user_id, username, first_name, last_name, joined_at, last_check)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := m.db.Exec(query,
		subscriber.UserID,
		subscriber.Username,
		subscriber.FirstName,
		subscriber.LastName,
		subscriber.JoinedAt.Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)
	return err
}

func (m *Monitor) deleteSubscriber(userID int64) error {
	_, err := m.db.Exec("DELETE FROM subscribers WHERE user_id = ?", userID)
	return err
}

func (m *Monitor) subscriberExists(userID int64) (bool, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM subscribers WHERE user_id = ?", userID).Scan(&count)
	return count > 0, err
}

func (m *Monitor) getSubscribersCount() (int, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM subscribers").Scan(&count)
	return count, err
}
