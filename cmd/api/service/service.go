package service

import (
	"encoding/base64"
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
	GetUser(string) (*storage.User, error)
	InsertRefreshToken(*storage.User) error
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

func (s *Service) Authenticate(uuid, addr string) (*Tokens, error) {
	user, err := s.Storage.GetUser(uuid)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return nil, fmt.Errorf("user not found")
	}

	TTLAccess, err := time.ParseDuration(os.Getenv("TTL_ACCESS"))
	if err != nil {
		return nil, err
	}

	payload := jwt.MapClaims{
		"uuid": uuid,
		"addr": addr,
		"exp":  time.Now().Add(time.Hour * TTLAccess),
	}

	tokenJWTString, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payload).SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("problem with signing token: ", err)
	}

	refreshToken := []byte(generateRefreshToken())

	refreshTokenBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(refreshToken)))
	base64.StdEncoding.Encode(refreshTokenBase64, refreshToken)

	tokens := &Tokens{
		AccessToken:  tokenJWTString,
		RefreshToken: refreshTokenBase64,
	}

	TTLRefresh, err := time.ParseDuration(os.Getenv("TTL_REFRESH"))
	if err != nil {
		return nil, err
	}

	brcyptRefreshToken, err := bcrypt.GenerateFromPassword(refreshTokenBase64, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.ExpiresAtRefresh = time.Now().Add(time.Hour * TTLRefresh)
	user.RefreshToken = brcyptRefreshToken

	err = s.Storage.InsertRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func (s *Service) Refresh() {

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
