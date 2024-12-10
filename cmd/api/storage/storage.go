package storage

import (
	"time"

	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

type User struct {
	Uuid             string
	Email            string
	RefreshToken     []byte
	ExpiresAtRefresh time.Time
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) GetUser(uuid string) (*User, error) {
	var user *User
	res := s.DB.First(user, "uuid = ?", uuid)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) InsertRefreshToken(user *User) error {
	return s.DB.Updates(user).Error
}
