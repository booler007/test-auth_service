package service

import (
	"encoding/base64"
	"fmt"
	"log"
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

type Emailer interface {
	NotificationNewIP(string) error
}

type Service struct {
	Storage Storager
	Email   Emailer
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken []byte `json:"refresh_token"`
}

func NewService(str Storager, eml Emailer) *Service {
	return &Service{str, eml}
}

func (s *Service) Authenticate(uuid, addrIP string) (*Tokens, error) {
	user, err := s.Storage.GetUserByUUID(uuid)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user not found")
	}

	return s.GenerateTokensAndSetSession(uuid, addrIP)
}

func (s *Service) RefreshTokens(refreshToken, addrIP string) (*Tokens, error) {
	hashRefreshToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	session, err := s.Storage.GetSessionByRefreshToken(hashRefreshToken)
	if err != nil {
		return nil, err
	}

	if session.IP != addrIP {
		user, err := s.Storage.GetUserByUUID(session.UserID)
		if err != nil {
			return nil, err
		}

		err = s.Email.NotificationNewIP(user.Email)
		if err != nil {
			log.Println(err)
		}
	}

	newTokens, err := s.GenerateTokensAndSetSession(session.UserID, session.IP)
	if err != nil {
		return nil, err
	}

	err = s.Storage.DeleteSession(session.ID)
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *Service) GenerateTokensAndSetSession(uuid, addrIP string) (*Tokens, error) {
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
		return nil, fmt.Errorf("problem with signing token: %s", err.Error())
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

	newSession := &storage.Session{
		UserID:       uuid,
		RefreshToken: bcryptRefreshToken,
		ExpiredAt:    time.Now().Add(time.Hour * TTLRefresh),
		IP:           addrIP,
	}

	err = s.Storage.SetSession(newSession)
	if err != nil {
		return nil, err
	}

	newTokens := &Tokens{
		AccessToken:  tokenJWTString,
		RefreshToken: refreshTokenBase64,
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
