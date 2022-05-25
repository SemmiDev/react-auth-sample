package main

import (
	"context"
	"database/sql"
)

type DataStore struct {
	db *sql.DB
}

func NewDataStore(db *sql.DB) *DataStore {
	return &DataStore{db}
}

const createAccount = `
INSERT INTO accounts (
  username,
  email,
  hashed_password,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
`

func (s *DataStore) CreateAccount(ctx context.Context, errChan chan<- error, arg *Account) {
	_, err := s.db.ExecContext(ctx, createAccount, arg.Username, arg.Email, arg.HashedPassword, arg.CreatedAt, arg.UpdatedAt)
	errChan <- err
}

const getAccount = `
SELECT username, email, hashed_password, is_email_verified, is_deleted, created_at, updated_at FROM accounts
WHERE username = $1 AND is_deleted = false LIMIT 1
`

func (s *DataStore) GetAccountByUsername(ctx context.Context, errChan chan<- error, account chan<- *Account, username string) {
	row := s.db.QueryRowContext(ctx, getAccount, username)
	var i Account
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashedPassword,
		&i.IsEmailVerified,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	account <- &i
	errChan <- err
}
