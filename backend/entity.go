package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashed_password"`

	IsEmailVerified bool  `json:"is_email_verified"`
	IsDeleted       bool  `json:"is_active"`
	CreatedAt       int64 `json:"created_at"`
	UpdatedAt       int64 `json:"updated_at"`
	DeletedAt       int64 `json:"deleted_at"`
}

func NewAccount(username, email, password string) (*Account, error) {
	account := &Account{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	if err := account.HashPassword(password); err != nil {
		return nil, err
	}
	return account, nil
}

func (a *Account) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	a.HashedPassword = hashedPassword
	return nil
}

func (a *Account) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword(a.HashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}
