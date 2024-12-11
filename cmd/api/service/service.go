package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"authentication_medods/cmd/api/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Storager interface {
	GetUserByUUID(string) (*storage.User, error)
	GetSessionByRefreshToken([]byte) (*storage.Session, error)
	SetSession(session *storage.Session) error
	DeleteSession(int) error
}

type Service struct {
	Storage Storager
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken []byte `json:"refresh_token"`
}

func NewService(str Storager) *Service {
	return &Service{str}
}

func (s *Service) Authenticate(uuid, addrIP string) (*Tokens, error) {
	user, err := s.Storage.GetUserByUUID(uuid)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user not found")
	}

	return s.generateTokens(uuid, addrIP)
}

func (s *Service) RefreshTokens(refreshToken, addrIP string) (*Tokens, error) {
	hashRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	session, err := s.Storage.GetSessionByRefreshToken(hashRefreshToken)
	if err != nil {
		return nil, err
	}

	if session.IP != addrIP {
		//TODO: отправка письма юзеру
		return nil, errors.New("new IP address, please confirm it")
	}

	newTokens, err := s.generateTokens(session.UserID, session.IP)
	if err != nil {
		return nil, err
	}

	err = s.Storage.DeleteSession(session.ID)
	if err != nil {
		return nil, err
	}

	return newTokens, nil

}

func generateRefreshToken() string {
	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789")

	var b strings.Builder
	for i := 0; i < 20; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func (s *Service) generateTokens(uuid, addrIP string) (*Tokens, error) {
	TTLAccess, err := time.ParseDuration(os.Getenv("TTL_ACCESS"))
	if err != nil {
		return nil, err
	}

	payload := jwt.MapClaims{
		"uuid":   uuid,
		"addrIP": addrIP,
		"exp":    time.Now().Add(time.Hour * TTLAccess),
	}

	tokenJWTString, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payload).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("problem with signing token: ", err)
	}

	refreshToken := []byte(generateRefreshToken())

	refreshTokenBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(refreshToken)))
	base64.StdEncoding.Encode(refreshTokenBase64, refreshToken)

	TTLRefresh, err := time.ParseDuration(os.Getenv("TTL_REFRESH"))
	if err != nil {
		return nil, err
	}

	bcryptRefreshToken, err := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	session := &storage.Session{
		UserID:       uuid,
		RefreshToken: bcryptRefreshToken,
		ExpiredAt:    time.Now().Add(time.Hour * TTLRefresh),
		IP:           addrIP,
	}

	err = s.Storage.SetSession(session)
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{
		AccessToken:  tokenJWTString,
		RefreshToken: refreshTokenBase64,
	}

	return tokens, nil
}
