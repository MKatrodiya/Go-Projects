package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func CreatePostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=meet sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		account_number int,
		balance int,
		created_at timestamp
	)`

	_, err := s.db.Query(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `insert into account 
	(first_name, last_name, account_number, balance, created_at)
	values 
	($1, $2, $3, $4, $5)`

	res, err := s.db.Query(query,
		account.FirstName,
		account.LastName,
		account.AccountNumber,
		account.Balance,
		account.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", res)

	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `select * from account`

	res, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for res.Next() {
		account := new(Account)

		err := res.Scan(&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.AccountNumber,
			&account.Balance,
			&account.CreatedAt)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
