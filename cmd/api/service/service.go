package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"authentication_medods/cmd/api/storage"

	"github.com/golang-jwt/jwt/v5"
)

type Storager interface {
	GetUserByUUID(string) (*storage.User, error)
	AddUsedRefreshToken(string) error
	IsRefreshTokenValid(string) bool
}

type Emailer interface {
	NotificationNewIP(string) error
}

type Service struct {
	Storage    Storager
	Email      Emailer
	TTLAccess  time.Duration
	TTLRefresh time.Duration
	JWTSecret  []byte
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(str Storager, eml Emailer, acs, rfr time.Duration, srt []byte) *Service {
	return &Service{
		str,
		eml,
		acs,
		rfr,
		srt,
	}
}

func (s *Service) Authenticate(uuid, addrIP string) (*Tokens, error) {
	if _, err := s.Storage.GetUserByUUID(uuid); err != nil {
		return nil, err
	}

	return generateTokens(uuid, addrIP, s.TTLAccess, s.TTLRefresh, s.JWTSecret)
}

func (s *Service) RefreshTokens(refreshToken, addrIP string) (*Tokens, error) {
	decodeRefreshTokenBase64, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(string(decodeRefreshTokenBase64), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.JWTSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	sign := strings.Split(string(decodeRefreshTokenBase64), ".")[2]
	if !s.Storage.IsRefreshTokenValid(sign) {
		return nil, errors.New("refresh token is invalid")
	}

	if claims["addrIP"] != addrIP {
		user, err := s.Storage.GetUserByUUID(claims["uuid"].(string))
		if err != nil {
			return nil, err
		}

		err = s.Email.NotificationNewIP(user.Email)
		if err != nil {
			log.Println(err)
		}
	}

	newTokens, err := generateTokens(claims["uuid"].(string), addrIP, s.TTLAccess, s.TTLRefresh, s.JWTSecret)
	if err != nil {
		return nil, err
	}

	if err = s.Storage.AddUsedRefreshToken(sign); err != nil {
		return nil, err
	}

	return newTokens, nil
}

func generateTokens(uuid, addrIP string, ttlAccess, ttlRefresh time.Duration, jwtSecret []byte) (*Tokens, error) {
	payloadAccess := jwt.MapClaims{
		"uuid":      uuid,
		"addrIP":    addrIP,
		"expiredAt": time.Now().Add(ttlAccess).Format(time.DateTime),
	}

	payloadRefresh := jwt.MapClaims{
		"uuid":      uuid,
		"addrIP":    addrIP,
		"expiredAt": time.Now().Add(ttlRefresh).Format(time.DateTime),
	}

	tokenAccess, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payloadAccess).SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("problem with signing access token: %s", err.Error())
	}

	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payloadRefresh).SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("problem with signing access token: %s", err.Error())
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(tokenRefresh))

	newTokens := &Tokens{
		AccessToken:  tokenAccess,
		RefreshToken: refreshTokenBase64,
	}

	return newTokens, nil
}
