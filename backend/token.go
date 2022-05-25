package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

func EpochToTimeZone(epoch int64, timeZone string) (time.Time, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(epoch, 0).In(loc), nil
}

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  int64     `json:"issued_at"`
	ExpiredAt int64     `json:"expired_at"`
	TimeZone  string    `json:"time_zone"`
}

func NewPayload(username string, timeZone string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  now.UTC().Unix(),
		ExpiredAt: now.Add(duration).UTC().Unix(),
		TimeZone:  timeZone,
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	now, err := EpochToTimeZone(time.Now().Unix(), payload.TimeZone)
	if err != nil {
		return err
	}

	expiredAt, err := EpochToTimeZone(payload.ExpiredAt, payload.TimeZone)
	if err != nil {
		return err
	}

	if now.After(expiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type Maker interface {
	CreateToken(username string, timeZone string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, timeZone string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, timeZone, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
