package storage

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string
	Name         string

	// Telegram linking
	TelegramUsername string `gorm:"uniqueIndex"`

	IsAdmin bool
}

func (db *DB) CreateUser(user *User) error {
	return db.conn.Create(user).Error
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	var user User
	err := db.conn.Where("email = ?", email).Find(&user).Error
	return user, err
}

func (db * DB) GetUserByTelegram(tgUsername string) (User, error) {
	var user User
	err := db.conn.Where("telegram_username = ?", tgUsername).Find(&user).Error
	return user, err
}
