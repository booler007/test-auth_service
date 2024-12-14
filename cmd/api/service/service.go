package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"authentication_medods/cmd/api/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Storager interface {
	GetUserByUUID(string) (*storage.User, error)
	AddUsedRefreshToken(storage.UsedRefreshTokens) error
	GetRefreshTokensByTime(time.Time) ([]storage.UsedRefreshTokens, error)
}

type Emailer interface {
	SendNotificationNewIP(string) error
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

var (
	ErrInvalidRefreshToken      = errors.New("refresh token is invalid")
	ErrProblemWithSigningTokens = errors.New("problem with signing token")
)

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

	createdTime, err := time.Parse(time.RFC3339Nano, claims["createdAt"].(string))
	if err != nil {
		return nil, err
	}

	fingerprintsOfRT, err := s.Storage.GetRefreshTokensByTime(createdTime)
	if err != nil {
		return nil, err
	}
	for _, fingerprint := range fingerprintsOfRT {
		if err = bcrypt.CompareHashAndPassword(fingerprint.Hash, decodeRefreshTokenBase64[len(decodeRefreshTokenBase64)-49:]); err == nil {
			return nil, ErrInvalidRefreshToken
		}
	}

	hashedJWT, err := bcrypt.GenerateFromPassword(decodeRefreshTokenBase64[len(decodeRefreshTokenBase64)-49:], bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	fingerprintOfUsedRT := storage.UsedRefreshTokens{
		CreatedAt: createdTime,
		Hash:      hashedJWT,
	}

	if err = s.Storage.AddUsedRefreshToken(fingerprintOfUsedRT); err != nil {
		return nil, err
	}

	if claims["addrIP"] != addrIP {
		user, err := s.Storage.GetUserByUUID(claims["uuid"].(string))
		if err != nil {
			return nil, err
		}

		err = s.Email.SendNotificationNewIP(user.Email)
		if err != nil {
			log.Println(err)
		}
	}

	newTokens, err := generateTokens(claims["uuid"].(string), addrIP, s.TTLAccess, s.TTLRefresh, s.JWTSecret)
	if err != nil {
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
		"createdAt": time.Now(),
	}

	tokenAccess, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payloadAccess).SignedString(jwtSecret)
	if err != nil {
		return nil, ErrProblemWithSigningTokens
	}

	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodHS512, payloadRefresh).SignedString(jwtSecret)
	if err != nil {
		return nil, ErrProblemWithSigningTokens
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(tokenRefresh))

	newTokens := &Tokens{
		AccessToken:  tokenAccess,
		RefreshToken: refreshTokenBase64,
	}

	return newTokens, nil
}
