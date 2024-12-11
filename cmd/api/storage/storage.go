package storage

import (
	"fmt"
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
	ExpiredAtRefresh time.Time
}

type Session struct {
	ID           int
	UserID       string
	RefreshToken []byte
	ExpiredAt    time.Time
	IP           string
}

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) GetUserByUUID(uuid string) (*User, error) {
	var user *User
	res := s.DB.First(user, "uuid = ?", uuid)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (s *Storage) GetSessionByRefreshToken(refreshToken []byte) (*Session, error) {
	var session *Session
	res := s.DB.First(session, "refresh_token = ?", refreshToken)
	if res.Error != nil {
		return nil, res.Error
	}

	if session.ExpiredAt.After(time.Now()) {
		return session, nil
	}

	return nil, fmt.Errorf("invalid refresh token")
}

func (s *Storage) SetSession(ss *Session) error {
	return s.DB.Create(ss).Error
}

func (s *Storage) DeleteSession(id int) error {
	return s.DB.Where("id = ?", id).Delete(&Session{}).Error
}
