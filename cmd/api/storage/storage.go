package storage

import (
	"errors"

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
	Signature string
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

func (s *Storage) AddUsedRefreshToken(sgnt string) error {
	return s.DB.Create(&UsedRefreshTokens{sgnt}).Error
}

func (s *Storage) IsRefreshTokenValid(signature string) bool {
	return s.DB.First(&UsedRefreshTokens{}, "signature = ?", signature).RowsAffected == 0
}
