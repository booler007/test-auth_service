package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

type User struct {
	Uuid  string
	Email string
}

type UsedRefreshTokens struct {
	CreatedAt time.Time
	Hash      []byte
}

var ErrUserNotFound = errors.New("user not found")

func NewStorage(db *gorm.DB) *Storage {
	return &Storage{db}
}

func (s *Storage) GetUserByUUID(uuid string) (*User, error) {
	var user User
	res := s.DB.First(&user, "uuid = ?", uuid)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, res.Error
	}

	return &user, nil
}

func (s *Storage) AddUsedRefreshToken(urt UsedRefreshTokens) error {
	return s.DB.Create(&urt).Error
}

func (s *Storage) GetRefreshTokensByTime(time time.Time) ([]UsedRefreshTokens, error) {
	var hashes []UsedRefreshTokens
	res := s.DB.Where("created_at = ?", time).Find(&hashes)
	if res.Error != nil {
		return nil, res.Error
	}

	return hashes, nil
}
