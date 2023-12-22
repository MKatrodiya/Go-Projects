package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() error
}

type PostgresStore struct {
	db *sql.DB
}

func CreatePostgresStore() (*PostgresStore, error) {

}
